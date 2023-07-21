// formatTimestamp formats the timestamp in locale string
export const formatTimestamp = function (timestamp) {
  return new Date(timestamp).toLocaleString();
};

export const getTimezoneName = function () {
  return Intl.DateTimeFormat().resolvedOptions().timeZone;
};

export const formatDayForInput = function (date) {
  let month, day, year;
  (month = '' + (date.getMonth() + 1)),
    (day = '' + date.getDate()),
    (year = date.getFullYear());

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

export const addTimeLeadZero = function (time) {
  return ('0' + time).slice(-2);
};

export type TimeBetweenUnits = {
  seconds: number;
  minutes: number;
  hours: number;
  days: number;
};

export const timeUnitsBetween = function (
  startDate,
  endDate,
): TimeBetweenUnits {
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
