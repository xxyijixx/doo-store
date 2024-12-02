import http from '../index'
import * as common from '@/api/interface/common'


export const getDetail = (key: string) =>  {
    return http.get<common.Detafocus>(`/api/v1/apps/${key}/detail`)
}

export const postInstall = (key: string, params: object) =>  {
    return http.post<common.Detail>(`/api/v1/apps/${key}`, params)
}

export const putAppStatus = (key: string, params: object) => {
    return http.put(`/api/v1/apps/${key}`, params)
}

export const getInsParams = (id: string|number) =>  {
    return http.get<common.getEdit>(`/api/v1/apps/installed/${id}/params`)
}

export const putInsParams = (id: string|number, params: object) =>  {
    return http.put<common.getEdit>(`/api/v1/apps/installed/${id}/params`, params)
}

export const getTags = () => {
    return http.get<common.Tag[]>(`/api/v1/apps/tags`)
}

export const getLogs = (id: string|number, params: object) =>  {
    return http.get<string>(`/api/v1/apps/installed/${id}/logs`,params)
}

export const deleteApp = (key: string, params?: object) =>  {
    return http.delete(`/api/v1/apps/${key}`, params)
}

export const getAppList = (params: object) =>  {
    return http.get<common.Data<common.Item>>(`/api/v1/apps`,  params )
}

export const getInstalledAppList = (params: object) =>  {
    return http.get<common.Data<common.Item>>(`/api/v1/apps/installed`,  params )
}


