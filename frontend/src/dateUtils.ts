// formatTimestamp formats the timestamp in locale string
export const formatTimestamp = function (timestamp) {
    return new Date(timestamp).toLocaleString()
}

export const getTimezoneName = function () {
    return Intl.DateTimeFormat().resolvedOptions().timeZone
}

export const formatDayForInput = function (date) {
    let month, day, year
    ;(month = '' + (date.getMonth() + 1)),
        (day = '' + date.getDate()),
        (year = date.getFullYear())

    if (month.length < 2) month = '0' + month
    if (day.length < 2) day = '0' + day

    return [year, month, day].join('-')
}

export const subtractDays = function (date, days) {
    return new Date(
        date.getFullYear(),
        date.getMonth(),
        date.getDate() - days,
        date.getHours(),
        date.getMinutes(),
        date.getSeconds(),
        date.getMilliseconds(),
    )
}
