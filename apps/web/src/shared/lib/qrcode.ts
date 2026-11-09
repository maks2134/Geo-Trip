/**
 * Self-contained SVG QR Code Generator (Byte Mode, Version 1-4, Low Error Correction)
 * Zero external dependencies, pure TypeScript, rendering crisp scalable SVG.
 */

// GF(256) math tables for Reed-Solomon encoding
const GF256_EXP = new Uint8Array(512);
const GF256_LOG = new Uint8Array(256);

(() => {
  let x = 1;
  for (let i = 0; i < 255; i++) {
    GF256_EXP[i] = x;
    GF256_LOG[x] = i;
    x <<= 1;
    if (x & 256) x ^= 0x11d;
  }
  for (let i = 255; i < 512; i++) {
    GF256_EXP[i] = GF256_EXP[i - 255];
  }
})();

function gfMul(x: number, y: number): number {
  if (x === 0 || y === 0) return 0;
  return GF256_EXP[GF256_LOG[x] + GF256_LOG[y]];
}

function rsGenPoly(numEcCodewords: number): Uint8Array {
  let poly = new Uint8Array([1]);
  for (let i = 0; i < numEcCodewords; i++) {
    const next = new Uint8Array(poly.length + 1);
    const root = GF256_EXP[i];
    for (let j = 0; j < poly.length; j++) {
      next[j] ^= gfMul(poly[j], root);
      next[j + 1] ^= poly[j];
    }
    poly = next;
  }
  return poly;
}

function rsCalculateEC(data: Uint8Array, numEcCodewords: number): Uint8Array {
  const gen = rsGenPoly(numEcCodewords);
  const out = new Uint8Array(numEcCodewords);
  for (let i = 0; i < data.length; i++) {
    const feedback = data[i] ^ out[0];
    for (let j = 0; j < numEcCodewords - 1; j++) {
      out[j] = out[j + 1] ^ gfMul(gen[gen.length - 2 - j], feedback);
    }
    out[numEcCodewords - 1] = gfMul(gen[0], feedback);
  }
  return out;
}

// Version table: [version, totalCodewords, dataCodewords, ecCodewords]
const VERSIONS = [
  { ver: 1, size: 21, total: 26, data: 19, ec: 7 },
  { ver: 2, size: 25, total: 44, data: 34, ec: 10 },
  { ver: 3, size: 29, total: 70, data: 55, ec: 15 },
  { ver: 4, size: 33, total: 100, data: 80, ec: 20 },
  { ver: 5, size: 37, total: 134, data: 108, ec: 26 },
];

export function generateQRCodeMatrix(text: string): boolean[][] {
  const utf8 = new TextEncoder().encode(text);
  let vInfo = VERSIONS[0];
  for (const v of VERSIONS) {
    if (v.data - 2 >= utf8.length) {
      vInfo = v;
      break;
    }
  }

  // 1. Bitstream encoding: Byte mode (0100) + 8-bit length + data + terminator
  const bits: number[] = [0, 1, 0, 0];
  const charCount = utf8.length;
  for (let i = 7; i >= 0; i--) bits.push((charCount >> i) & 1);
  for (let i = 0; i < utf8.length; i++) {
    const b = utf8[i];
    for (let j = 7; j >= 0; j--) bits.push((b >> j) & 1);
  }
  // Terminator
  while (bits.length % 8 !== 0) bits.push(0);

  // Convert bits to data codewords
  const dataCodewords = new Uint8Array(vInfo.data);
  let byteIndex = 0;
  for (let i = 0; i < bits.length; i += 8) {
    let b = 0;
    for (let j = 0; j < 8; j++) {
      b = (b << 1) | (bits[i + j] || 0);
    }
    if (byteIndex < dataCodewords.length) {
      dataCodewords[byteIndex++] = b;
    }
  }

  // Padding
  const pad = [0xec, 0x11];
  let padIdx = 0;
  while (byteIndex < dataCodewords.length) {
    dataCodewords[byteIndex++] = pad[padIdx++ % 2];
  }

  // 2. Error correction
  const ecCodewords = rsCalculateEC(dataCodewords, vInfo.ec);

  // 3. Assemble full codeword stream
  const fullCodewords = new Uint8Array(vInfo.total);
  fullCodewords.set(dataCodewords, 0);
  fullCodewords.set(ecCodewords, dataCodewords.length);

  // 4. Matrix creation
  const size = vInfo.size;
  const matrix: (boolean | null)[][] = Array.from({ length: size }, () =>
    Array(size).fill(null)
  );

  // Draw finder patterns (7x7) at 3 corners
  const drawFinder = (r: number, c: number) => {
    for (let i = -1; i <= 7; i++) {
      for (let j = -1; j <= 7; j++) {
        const nr = r + i;
        const nc = c + j;
        if (nr < 0 || nr >= size || nc < 0 || nc >= size) continue;
        if (i >= 0 && i <= 6 && j >= 0 && j <= 6) {
          if (i === 0 || i === 6 || j === 0 || j === 6 || (i >= 2 && i <= 4 && j >= 2 && j <= 4)) {
            matrix[nr][nc] = true;
          } else {
            matrix[nr][nc] = false;
          }
        } else {
          matrix[nr][nc] = false; // Separator
        }
      }
    }
  };

  drawFinder(0, 0);
  drawFinder(0, size - 7);
  drawFinder(size - 7, 0);

  // Timing patterns
  for (let i = 8; i < size - 8; i++) {
    matrix[6][i] = i % 2 === 0;
    matrix[i][6] = i % 2 === 0;
  }

  // Dark module
  matrix[4 * vInfo.ver + 9][8] = true;

  // Mask reserve: Format information area
  for (let i = 0; i < 9; i++) {
    if (matrix[8][i] === null) matrix[8][i] = false;
    if (matrix[i][8] === null) matrix[i][8] = false;
    if (matrix[8][size - 1 - i] === null) matrix[8][size - 1 - i] = false;
    if (matrix[size - 1 - i][8] === null) matrix[size - 1 - i][8] = false;
  }

  // 5. Place data codewords in 2-column zig-zag
  let bitIdx = 0;
  const totalBits = fullCodewords.length * 8;
  let dir = -1; // up
  let col = size - 1;
  let row = size - 1;

  while (col > 0) {
    if (col === 6) col--; // Skip vertical timing pattern
    for (let cOffset = 0; cOffset < 2; cOffset++) {
      const c = col - cOffset;
      if (matrix[row][c] === null) {
        let bit = false;
        if (bitIdx < totalBits) {
          const byte = fullCodewords[Math.floor(bitIdx / 8)];
          bit = ((byte >> (7 - (bitIdx % 8))) & 1) === 1;
          bitIdx++;
        }
        // Standard mask 0: (row + col) % 2 === 0
        const mask = (row + c) % 2 === 0;
        matrix[row][c] = mask ? !bit : bit;
      }
    }
    row += dir;
    if (row < 0 || row >= size) {
      dir = -dir;
      row += dir;
      col -= 2;
    }
  }

  // Final conversion of any remaining nulls to false
  return matrix.map((row) => row.map((cell) => cell ?? false));
}

export function generateQRCodeSVG(
  text: string,
  pixelSize = 64,
  fgColor = "#000000",
  bgColor = "#ffffff"
): string {
  const matrix = generateQRCodeMatrix(text);
  const size = matrix.length;
  let rects = "";

  for (let r = 0; r < size; r++) {
    for (let c = 0; c < size; c++) {
      if (matrix[r][c]) {
        rects += `<rect x="${c}" y="${r}" width="1" height="1" fill="${fgColor}"/>`;
      }
    }
  }

  return `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 ${size} ${size}" width="${pixelSize}" height="${pixelSize}" style="shape-rendering: crispEdges;">
    <rect width="${size}" height="${size}" fill="${bgColor}"/>
    ${rects}
  </svg>`;
}
