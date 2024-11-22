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
}

export interface Detail {
    label: string
    env_key: string
    values: object
    type: string
    rule: string
    required: boolean
}

export interface Detafocus {
    id: number
    created_at: string
    updated_at: string
    app_id: number
    repo: string
    version: string
    depends_version: string
    docker_compose: string
    nginx_config: string
    status: string
    params: Params
  }
  
  export interface Params {
    form_fields: FormField[]
  }

export interface FormField {
    default: string
    label: string
    env_key: string
    key: string
    value: string
    values: object
    type: string
    rule: string
    required: boolean
}

export interface getEdit {
    params: editForm[]
    docker_compose: string
    cpus: string
    memory_limit: string
  }

  export interface editForm {
    default: string
    label: string
    env_key: string
    key: string
    value: string
    values: object
    type: string
    rule: string
    required: boolean
  }
  