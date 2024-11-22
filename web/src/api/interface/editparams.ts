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

