export default {
  stringifyQuery: (query) => {
    const keys = Object.keys(query)
    return keys.map(key => `${key}=${encodeURIComponent(query[key])}`).join('&')
  }
}
