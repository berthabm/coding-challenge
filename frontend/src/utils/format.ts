/**
 * Formats a number for display:
 * - Integers → no decimals   (5.0 → "5")
 * - Floats   → max 3 decimals, trailing zeros stripped  (-0.1231 → "-0.123", 2.5 → "2.5")
 */
export function formatNumber(v: number): string {
  if (Number.isInteger(v)) return v.toString();
  // toFixed(3) then strip trailing zeros
  return parseFloat(v.toFixed(3)).toString();
}
