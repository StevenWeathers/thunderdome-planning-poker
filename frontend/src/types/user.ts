export type User = {
    avatar: string
    company: string
    country: string
    createdDate: string
    disabled: boolean
    email: string
    gravatarHash: string
    id: string
    jobTitle: string
    lastActive: string
    locale: string
    mfaEnabled: boolean
    name: string
    notificationsEnabled: boolean
    rank: string
    updatedDate: string
    verified: boolean
}

export type UserAPIKey = {
    active: boolean
    apiKey: string
    createdDate: string
    id: string
    name: string
    prefix: string
    updatedDate: string
    userEmail: string
    userId: string
    userName: string
}
