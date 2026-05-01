// formatTimestamp formats the timestamp in locale string
export const formatTimestamp = function (timestamp) {
  return new Date(timestamp).toLocaleString();
};

export const getTimezoneName = function () {
  return Intl.DateTimeFormat().resolvedOptions().timeZone || 'America/New_York'; // add fallback for users whose timezone is undefined
};

type DatePartsForInput = {
  year: string;
  month: string;
  day: string;
};

const getDatePartsForTimezone = (date: Date, timeZone: string): DatePartsForInput => {
  const formatter = new Intl.DateTimeFormat('en-CA', {
    timeZone,
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
  });

  return formatter.formatToParts(date).reduce(
    (parts, part) => {
      if (part.type === 'year' || part.type === 'month' || part.type === 'day') {
        parts[part.type] = part.value;
      }

      return parts;
    },
    { year: '', month: '', day: '' },
  );
};

export const formatDayForInput = function (date: Date, timeZone?: string): string {
  if (timeZone) {
    const { year, month, day } = getDatePartsForTimezone(date, timeZone);

    return [year, month, day].join('-');
  }

  let month, day, year;
  ((month = '' + (date.getMonth() + 1)), (day = '' + date.getDate()), (year = date.getFullYear()));

  if (month.length < 2) month = '0' + month;
  if (day.length < 2) day = '0' + day;

  return [year, month, day].join('-');
};

export const subtractDays = function (date, days) {
  return new Date(
    date.getFullYear(),
    date.getMonth(),
    date.getDate() - days,
    date.getHours(),
    date.getMinutes(),
    date.getSeconds(),
    date.getMilliseconds(),
  );
};

export const addMinutesToDate = (date, n) => {
  const d = new Date(date);
  d.setTime(d.getTime() + n * 60000);
  return d;
};

export const addTimeLeadZero = function (time) {
  return ('0' + time).slice(-2);
};

export type TimeBetweenUnits = {
  seconds: number;
  minutes: number;
  hours: number;
  days: number;
};

export const timeUnitsBetween = function (startDate, endDate): TimeBetweenUnits {
  let delta = Math.abs(endDate - startDate) / 1000;
  const unitDivisions: Array<{
    key: string;
    value: number;
  }> = [
    { key: 'days', value: 24 * 60 * 60 },
    { key: 'hours', value: 60 * 60 },
    { key: 'minutes', value: 60 },
    { key: 'seconds', value: 1 },
  ];
  return unitDivisions.reduce(
    (acc: TimeBetweenUnits, arrValue) => (
      (acc[arrValue.key] = Math.floor(delta / arrValue.value)),
      (delta -= acc[arrValue.key] * arrValue.value),
      acc
    ),
    {
      seconds: 0,
      minutes: 0,
      hours: 0,
      days: 0,
    },
  );
};
