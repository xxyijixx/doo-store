export interface Root<T> {
    code: number
    msg: string
    data: T
}

export interface Data<T> {
    total: number
    items: T[]
}

export interface Item {
    id: number
    created_at: string
    updated_at: string
    name: string
    key: string
    icon: string
    description: string
    github: string
    class: string
    depends_version: string
    sort: number
    status: string
    message: string
}

export interface form_fields<T> {
    label: string
    env_key: string
    values: T
    type: string
    rule: string
    required: boolean
}

export interface Params<T> {
    default: string
    label: string
    env_key: string
    key: string
    value: string
    values: T
    type: string
    rule: string
    required: boolean
}