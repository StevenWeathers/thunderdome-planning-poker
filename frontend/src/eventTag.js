// we don't want timeouts to google analytics to hold up page routing
function createFunctionWithTimeout(callback, opt_timeout) {
    var called = false
    function fn() {
        if (!called) {
            called = true
            callback()
        }
    }
    setTimeout(fn, opt_timeout || 250)
    return fn
}

export default function (action, category, label, cb = function () {}) {
    // provide fallback should gtag not be available
    const t =
        window.gtag ||
        function (evt, action, opts) {
            cb()
        }

    t('event', action, {
        event_category: category,
        event_label: label,
        event_callback: createFunctionWithTimeout(cb),
    })
}
