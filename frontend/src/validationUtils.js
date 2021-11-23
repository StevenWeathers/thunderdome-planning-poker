export const nameMin = 1
export const nameMax = 64
export const passMin = 6
export const passMax = 72
export const nameLenError = `Name must be between ${nameMin} and ${nameMax} characters.`
export const passLenError = `Password must be between ${passMin} and ${passMax} characters.`
export const passNotMatchError = `Password and Confirm Password do not match.`

export const validatePasswords = (password1, password2) => {
    let valid = true
    let error = ''

    if (password1.length < passMin || password1.length > passMax) {
        valid = false
        error = passLenError
    }

    if (password1 !== password2) {
        valid = false
        error = passNotMatchError
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
        error = nameLenError
    }

    return {
        valid,
        error,
    }
}

export const validateUserIsAdmin = user => {
    return user.rank === 'ADMIN' || user.rank === 'GENERAL'
}

export const validateUserIsRegistered = user => {
    return user.rank !== 'GUEST' && user.rank !== 'PRIVATE'
}
