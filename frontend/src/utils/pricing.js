// Central place for prices (RSD) used by the UI.
//
// IMPORTANT: Backend DB uses SR type names (PASOS, LICNA_KARTA),
// while some older code used EN names (PASSPORT, ID_CARD).
// We support BOTH so the UI never shows 0 by mistake.

export const PRICES = {
  // ID card
  LICNA_KARTA: 3500,
  ID_CARD: 3500,

  // Passport
  PASOS: 4200,
  PASSPORT: 4200,

  // Driver license (if you use it)
  VOZACKA: 4000,
  DRIVER_LICENSE: 4000,
};

export function priceFor(type) {
  if (!type) return 0;
  const key = String(type).trim().toUpperCase();
  return PRICES[key] ?? 0;
}

export function fmtRSD(n) {
  const v = Number(n ?? 0);
  if (Number.isNaN(v)) return "";
  return v.toLocaleString("sr-RS") + " RSD";
}
