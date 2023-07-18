export type Team = {
    createdDate: string
    id: string
    name: string
    updatedDate: string
}

export type TeamUser = {
    avatar: string
    email: string
    gravatarHash: string
    id: string
    name: string
    role: string
}

export type TeamCheckin = {
    blockers: string
    comments: Array<CheckinComment>
    createdDate: string
    discuss: string
    goalsMet: boolean
    id: string
    today: string
    updatedDate: string
    user: Array<TeamUser>
    yesterday: string
}

export type CheckinComment = {
    checkin_id: string
    comment: string
    created_date: string
    id: string
    updated_date: string
    user_id: string
}
