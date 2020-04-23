/**
 * Extends fetch with common inputs e.g. credentials, content-type
 * and checks response status/ok for common errors
 * @param {function} handle401 401 handler, e.g. redirect to login
 */
export default function(handle401) {
    /**
     * Wrapper around fetch
     * @param {string} endpoint the endpoint to fetch
     * @param {object} config the optional fetch config e.g. body for post
     */
    return function(endpoint, { body, ...customConfig } = {}) {
        const headers = { 'content-type': 'application/json' }

        const config = {
            method: body ? 'POST' : 'GET',
            credentials: 'same-origin',
            ...customConfig,
            headers: {
                ...headers,
                ...customConfig.headers,
            },
        }

        if (body) {
            config.body = JSON.stringify(body)
        }

        return fetch(endpoint, config).then(response => {
            if (response.status === 401) {
                handle401()
            }

            if (!response.ok) {
                throw Error(response.statusText)
            }

            return response
        })
    }
}
