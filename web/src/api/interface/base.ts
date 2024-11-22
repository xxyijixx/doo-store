/* eslint-disable @typescript-eslint/no-explicit-any */
// 通用返回
export interface Result {
    code: number
    msg: string,
}

// 通用返回数据
export interface ResultData<T = any> extends Result {
    [x: number]: string
    data?: T
}