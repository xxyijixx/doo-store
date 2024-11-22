/* eslint-disable @typescript-eslint/no-explicit-any */
/* eslint-disable @typescript-eslint/no-unused-vars */
"use client";

import { useForm } from "react-hook-form";
import { Button } from "@/components/ui/button";
import {
    Form,
    FormControl,
    FormItem,
    FormLabel,
    FormMessage,
} from "@/components/ui/form";
import { SheetClose } from "@/components/ui/sheet";
import { Input } from "@/components/ui/input";
import { Checkbox } from "@/components/ui/checkbox";
import { Label } from "@/components/ui/label";
import { HighConfig } from "@/components/drawer/highconfig";
import { Item } from "@/type.d/common";
import { useEffect, useState } from "react";
import * as http from "@/api/modules/fouceinter";
import { useTranslation } from "react-i18next";

interface ProfileFormProps {
    app: Item; // 接收 app 数据
    onInstallSuccess: () => void; // 新增回调
    onFalse: () => void; // 失败回调
}

interface Field {
    default: string;
    label: string;
    env_key: string;
    values: any;
    type: string;
    rule: string;
    required: boolean;
}

type FormValues = {
    [key: string]: string | number | boolean; // 根据实际情况调整类型
};

export function ProfileForm({
    app,
    onInstallSuccess,
    onFalse,
}: ProfileFormProps) {

    const { t } = useTranslation();
    const [dockerCompose, setDockerCompose] = useState<string>(""); // 存储docker_compose内容
    const [cpuLimit, setCpuLimit] = useState<string>("1"); // 默认值为 1
    const [memoryLimit, setMemoryLimit] = useState<string>("0"); // 默认值为 120M
    const [loading, setLoading] = useState<boolean>(false); // 加载状态
    const [error, setError] = useState<string>(""); // 错误信息
    // const [successMessage, setSuccessMessage] = useState<string>("");
    const [formFields, setFormFields] = useState<any[]>([]); // 存储 form_fields 数据

    const form = useForm<FormValues>({
        defaultValues: {},
    });

    const {
        setValue,
        handleSubmit,
        formState: { errors },
    } = form;

    // 发起请求获取 form_fields 内容
    useEffect(() => {
        if (app.key) {
            setLoading(true); // 开始加载
            setError(""); // 清空之前的错误
            // 发起 GET 请求
            http.getDetail(app.key)
                .then((response) => {
                    return response;
                })
                .then((data) => {
                    console.log("API Response:", data); // 调试：检查返回的数据
                    setDockerCompose(data.data?.docker_compose || "");

                    // 获取 form_fields 数组并存储
                    setFormFields(data.data?.params.form_fields || []);
                    // 设置 formFields 的默认值
                    data.data?.params.form_fields.forEach(
                        (field: Field) => {
                            const fieldName = field.env_key;
                            setValue(fieldName, field.default || ""); // 设置每个字段的默认值
                        }
                    );

                    setLoading(false); // 请求完成
                })
                .catch((_error) => {
                    setError(t("请求失败，请稍后重试")); // 错误处理
                    setLoading(false);
                });
        }
    }, [app.key, setValue]);

    // 点击重启按钮，进行安装操作
    const handleRestart = async () => {
        setLoading(true);
        setError(""); // 清除之前的错误信息
        // setSuccessMessage(""); // 清除之前的成功信息

        // 校验表单，确保字段不为空
        const isValid = await form.trigger(); // 校验所有表单字段
        if (!isValid) {
            console.log(t("表单校验失败！"));
            setError(t("请填完整的字段信息！")); // 提示用户填写完整字段
            setLoading(false);
            return;
        }
        const currentValues = form.getValues()
        console.log("表单数据:", currentValues)
        const params: { [key: string]: string | number | boolean } = {};
        formFields.forEach((field: Field) => {
            params[field.env_key] = currentValues[field.env_key]
        })
        try {
            // 构建请求体
            const requestBody = {
                cpus: cpuLimit,
                docker_compose: dockerCompose,
                memory_limit: memoryLimit,
                params: params, // 可根据需要填充参数
            };

            // 发送 POST 请求进行安装
            const response = await http.postInstall(app.key, requestBody);
            console.log("API Response:", response); // 调试：检查返回的数据
            if (response.code == 200) {
                console.log("安装成功！");
                onInstallSuccess(); // 成功时调用回调函数回到drawer组件
            } else {
                // 请求失败，显示错误消息
                console.log(error);
                // setError(response.message || "请求失败，请稍后重试");
                onFalse(); // 失败时调用回调函数回到drawer组件显示
            }
        } catch (error) {
            // 捕获网络错误
            setError(t("网络错误，请检查网络连接"));
            onFalse(); // 失败时调用回调函数回到drawer组件显示
        } finally {
            setLoading(false); // 请求完成，无论成功或失败都需要关闭加载状态
        }
    };

    return (
        <Form {...form} >
            <form className="space-y-8 " onSubmit={handleSubmit(handleRestart)}>
                {/* 动态渲染 form_fields */}
                {formFields.map((field, index) => {
                    // 如果 field 没有 name 属性，生成一个默认的 name
                    const fieldName = field.env_key

                    return (
                        <FormItem key={index}>
                            <FormLabel>{field.label}</FormLabel>
                            <FormControl>
                                <div className="w-full">
                                <Input
                                    id={fieldName}
                                    placeholder={t("请输入...")}
                                    {...form.register(fieldName, {
                                        required: t(`${field.label} 不能为空`),
                                        onChange: (_error) => {
                                            // 触发 onChange 时重新校验（即时校验）
                                            form.trigger(fieldName);
                                        },
                                        onBlur: () => {
                                            // 输入框失去焦点时触发校验
                                            form.trigger(fieldName);
                                        },
                                    })}
                                />
                                </div>
                            </FormControl>
                            <FormMessage>{errors[fieldName]?.message}</FormMessage>
                        </FormItem>
                    );
                })}

                <FormItem>
                    <div className="flex items-center space-x-2">
                        <Checkbox id="store" />
                        <Label htmlFor="store">{t('默认储存')}</Label>
                    </div>
                </FormItem>

                <FormItem>
                    <HighConfig
                        app={app}
                        dockerCompose={dockerCompose}
                        cpuLimit={cpuLimit} // 传递 cpuLimit
                        memoryLimit={memoryLimit} // 传递 memoryLimit
                        setDockerCompose={setDockerCompose}
                        setCpuLimit={setCpuLimit}
                        setMemoryLimit={setMemoryLimit}
                    />
                </FormItem>

                <div className="flex justify-start space-x-3">
                        <Button
                            type="submit"
                            variant="surely"
                            className="cursor-pointer"
                            onClick={handleRestart}
                            disabled={loading} //防止用户在请求期间重复点击按钮,没有这个条件影响局部加载，会全局加载覆盖出现问题
                        >
                            {t('安装')}
                        </Button>

                        <SheetClose className="cursor-pointer border border-input rounded-md bg-transparent text-sm text-gray-600 shadow-sm hover:bg-white hover:border-theme-color/85 hover:text-theme-color/85 h-9 px-5 py-2">
                            {t('取消')}
                        </SheetClose>

                    
                </div>
            </form>
        </Form>
    );
}
