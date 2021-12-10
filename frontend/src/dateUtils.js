// getTodaysDate gets today's date in YYYY-MM-DD format for the API input
export const getTodaysDate = function () {
    const date = new Date()
    return new Date(date.getTime() - date.getTimezoneOffset() * 60000)
        .toISOString()
        .substr(0, 10)
}

// formatTimestamp formats the timestamp in locale string
export const formatTimestamp = function (timestamp) {
    return new Date(timestamp).toLocaleString()
}

export const getTimezoneName = function () {
    return Intl.DateTimeFormat().resolvedOptions().timeZone
}
