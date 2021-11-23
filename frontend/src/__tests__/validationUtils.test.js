import {
    validatePasswords,
    validateName,
    validateUserIsAdmin,
    validateUserIsRegistered,
    nameMin,
    nameMax,
    passMin,
    passMax,
    nameLenError,
    passLenError,
    passNotMatchError,
} from '../validationUtils'

describe('Validation Utils', () => {
    describe('validateName', () => {
        it(`should allow name between ${nameMin} and ${nameMax} characters`, () => {
            expect(validateName('thor')).toEqual({
                valid: true,
                error: '',
            })
        })

        it(`should not allow name less than ${nameMin} characters`, () => {
            expect(validateName('')).toEqual({
                valid: false,
                error: nameLenError,
            })
        })

        it(`should not allow name greater than ${nameMax} characters`, () => {
            const reallyLongString =
                'kjadfjlkfsadjlkfdjklfsdjklasfjklfdsjlksdjklfdasjklfadsjkjklkjadfjlkfsadjlkfdjklfsdjklasfjklfdsjlksdjklfdasjklfadsjkjkl'
            expect(validateName(reallyLongString)).toEqual({
                valid: false,
                error: nameLenError,
            })
        })
    })

    describe('validatePassword', () => {
        it(`should allow name between ${passMax} and ${passMax} characters`, () => {
            expect(validatePasswords('password1', 'password1')).toEqual({
                valid: true,
                error: '',
            })
        })

        it(`should not allow name less than ${passMin} characters`, () => {
            expect(validatePasswords('tree', 'tree')).toEqual({
                valid: false,
                error: passLenError,
            })
        })

        it(`should not allow name greater than ${passMax} characters`, () => {
            const reallyLongString =
                'kjadfjdfjljskafdjkslajklfjklsdfljkadfslkfsadjlkfdjklfsdjklasfjklfdsjlksdjklfdasjklfadsjkjklkjadfjlkfsadjlkfdjklfsdjklasfjklfdsjlksdjklfdasjklfadsjkjkl'
            expect(
                validatePasswords(reallyLongString, reallyLongString),
            ).toEqual({
                valid: false,
                error: passLenError,
            })
        })

        it(`should not allow password1 and password2 miss matches`, () => {
            expect(
                validatePasswords('strongestavenger', 'lokiisajoke'),
            ).toEqual({
                valid: false,
                error: passNotMatchError,
            })
        })
    })

    describe('validateUserIsAdmin', () => {
        it('should return true for ADMIN rank user', () => {
            expect(validateUserIsAdmin({ rank: 'ADMIN' })).toEqual(true)
        })

        it('should return true for GENERAL rank user', () => {
            expect(validateUserIsAdmin({ rank: 'GENERAL' })).toEqual(true)
        })

        it('should return false for REGISTERED rank user', () => {
            expect(validateUserIsAdmin({ rank: 'REGISTERED' })).toEqual(false)
        })

        it('should return false for CORPORAL rank user', () => {
            expect(validateUserIsAdmin({ rank: 'CORPORAL' })).toEqual(false)
        })

        it('should return false for GUEST rank user', () => {
            expect(validateUserIsAdmin({ rank: 'GUEST' })).toEqual(false)
        })

        it('should return false for PRIVATE rank user', () => {
            expect(validateUserIsAdmin({ rank: 'PRIVATE' })).toEqual(false)
        })
    })

    describe('validateUserIsRegistered', () => {
        it('should return true for ADMIN rank user', () => {
            expect(validateUserIsRegistered({ rank: 'ADMIN' })).toEqual(true)
        })

        it('should return true for GENERAL rank user', () => {
            expect(validateUserIsRegistered({ rank: 'GENERAL' })).toEqual(true)
        })

        it('should return true for REGISTERED rank user', () => {
            expect(validateUserIsRegistered({ rank: 'REGISTERED' })).toEqual(
                true,
            )
        })

        it('should return true for CORPORAL rank user', () => {
            expect(validateUserIsRegistered({ rank: 'CORPORAL' })).toEqual(true)
        })

        it('should return false for GUEST rank user', () => {
            expect(validateUserIsRegistered({ rank: 'GUEST' })).toEqual(false)
        })

        it('should return false for PRIVATE rank user', () => {
            expect(validateUserIsRegistered({ rank: 'PRIVATE' })).toEqual(false)
        })
    })
})
