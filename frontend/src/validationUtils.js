const nameMin = 1
const nameMax = 64
const passMin = 6
const passMax = 72

export const validatePasswords = (password1, password2) => {
    let valid = true
    let error = ''

    if (password1.length < passMin || password1.length > passMax) {
        valid = false
        error = `Password must be between ${passMin} and ${passMax} characters.`
    }

    if (password1 !== password2) {
        valid = false
        error = `Password and Confirm Password do not match.`
    }

    return {
        valid,
        error,
    }
}

export const validateName = warriorName => {
    let valid = true
    let error = ''

    if (warriorName.length < nameMin || warriorName.length > nameMax) {
        valid = false
        error = `Name must be between ${nameMin} and ${nameMax} characters.`
    }

    return {
        valid,
        error,
    }
}
