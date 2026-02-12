export const convertToDate = (value: "5m" | "15m" | "1h" | "30h"): Date => {
  const now = new Date();

  if (value.endsWith("m")) {
    const minutes = parseInt(value);
    now.setMinutes(now.getMinutes() + minutes);
  }

  if (value.endsWith("h")) {
    const hours = parseInt(value);
    now.setHours(now.getHours() + hours);
  }

  return now;
};
