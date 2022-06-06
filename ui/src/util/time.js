import { DateTime } from "luxon";

export function formatIsoTimestamp(ts) {
  return DateTime.fromISO(ts)
    .setLocale("en-US")
    .toLocaleString(DateTime.DATETIME_SHORT);
}
