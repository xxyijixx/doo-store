/* eslint-disable react-hooks/exhaustive-deps */
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
import { FormField } from "@/api/interface/common"
import { useToast } from "@/hooks/use-toast";
import { FalseToaster } from "@/components/ui/toaster"

interface ProfileFormProps {
    app: Item; // 接收 app 数据
    onInstallSuccess: () => void; // 新增回调
    onFalse: () => void; // 失败回调
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
    const [formFields, setFormFields] = useState<FormField[]>([]); // 存储 form_fields 数据
    const { toast } = useToast(); // 引入 toast


    // 现有的状态和 hooks
    const [renderError, setRenderError] = useState<boolean>(false);

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
                    // 过滤无效的表单字段
                    const validFormFields = (response.data?.params.form_fields || [])
                        .filter(field => 
                            field.env_key && 
                            field.env_key.trim() !== '' && 
                            field.label && 
                            field.type
                        );

                    if (validFormFields.length === 0) {
                        // 如果没有有效字段，触发渲染错误
                        setRenderError(true);
                        // 添加 toast 弹窗提示
                        toast({
                            title: t("表单加载错误"),
                            description: t("未找到有效的表单字段，请检查应用配置"),
                            variant: "destructive",
                            duration: 5000, // 5秒后自动关闭
                        });
                        throw new Error("没有有效的表单字段");
                        
                    }

                    setDockerCompose(response.data?.docker_compose || "");
                    setFormFields(validFormFields);

                    // 仅为有效字段设置默认值
                    validFormFields.forEach((field) => {
                        const fieldName = field.env_key;
                        setValue(fieldName, field.default || "");
                    });

                    setLoading(false);
                    return response;
                })
                .then((data) => {
                    console.log("API Response:", data); // 调试：检查返回的数据
                    setDockerCompose(data.data?.docker_compose || "");

                    // 获取 form_fields 数组并存储
                    setFormFields(data.data?.params.form_fields || []);
                    // 设置 formFields 的默认值
                    data.data?.params.form_fields.forEach(
                        (field: FormField) => {
                            const fieldName = field.env_key;
                            setValue(fieldName, field.default || ""); // 设置每个字段的默认值
                        }
                    );
                    console.log("安装错误33333333333")
                    setLoading(false); // 请求完成
                })
                .catch((error) => {
                    console.error("getDetail 请求错误:", error);  // 添加详细的错误日志
                    setError(t("请求失败，请稍后重试")); // 错误处理
                    // 使用 toast 弹窗显示错误信息
                    toast({
                        title: t("请求失败"),
                        description: t("无法获取表单字段，请刷新后重试。"),
                        variant: "destructive", // 使用错误样式
                        duration: 5000, // 自动关闭时间，单位毫秒
                    });
                    setLoading(false);
                    
                });
        }
    }, [app.key, setValue]);

    // 点击按钮，进行安装操作
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
        formFields.forEach((field: FormField) => {
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
                console.log("安装错误111111111111")
                console.log(error);
                // setError(response.message || "请求失败，请稍后重试");
                onFalse(); // 失败时调用回调函数回到drawer组件显示
            }
        } catch (error) {
            // 捕获网络错误
            console.log("安装错误222222222222222")
            setError(t("网络错误，请检查网络连接"));
            onFalse(); // 失败时调用回调函数回到drawer组件显示
        } finally {
            setLoading(false); // 请求完成，无论成功或失败都需要关闭加载状态
        }
    };

    return (
        <>
        <FalseToaster />
        <Form {...form}>
            <form 
                className="space-y-8 relative overflow-visible lg:px-0 md:px-0 px-3 pb-3" 
                onSubmit={handleSubmit(handleRestart)}
                >
                {/* 动态渲染 form_fields */}
                {formFields.sort((a, b) => a.order - b.order).map((field, index) => {
                    const fieldName = field.env_key;

                    // 添加详细的字段验证和日志
                    console.group(`处理表单字段 ${index}`);
                    console.log('当前字段:', field);
                    
                    // 检查依赖关系
                    if (field.dependency) {
                        const dependentValue = form.watch(field.dependency.field);
                        if (dependentValue !== field.dependency.value) {
                            return null;
                        }
                    }

                    return (
                        
                        <FormItem key={index}>
                            <FormLabel>{field.label}</FormLabel>
                            <FormControl>
                                <>
                                    {field.type === 'select' && (
                                        <select
                                            className="flex h-9 w-full rounded-md border border-gray-200/60 bg-gray-200/60 px-3 py-1 text-sm shadow-sm transition-colors"
                                            id={fieldName}
                                            defaultValue={field.default}
                                            {...form.register(fieldName, {
                                                required: field.validation?.required && t(`${field.label} 不能为空`),
                                                onChange: () => form.trigger(fieldName),
                                                onBlur: () => form.trigger(fieldName),
                                            })}
                                        >
                                            {field.options?.map((option, i) => (
                                                <option key={i} value={option.value}>
                                                    {option.label}
                                                </option>
                                            ))}
                                        </select>
                                    )}
                                    {field.type === 'radio' && field.options && (
                                        <div className="flex gap-4">
                                            {field.options.map((option, i) => (
                                                <div key={i} className="flex items-center space-x-2">
                                                    <input
                                                        type="radio"
                                                        id={`${fieldName}-${i}`}
                                                        value={option.value}
                                                        {...form.register(fieldName, {
                                                            required: field.validation?.required && t(`${field.label} 不能为空`),
                                                        })}
                                                    />
                                                    <label htmlFor={`${fieldName}-${i}`}>{option.label}</label>
                                                </div>
                                            ))}
                                        </div>
                                    )}
                                    {field.type === 'checkbox' && (
                                        <Checkbox
                                            id={fieldName}
                                            {...form.register(fieldName)}
                                        />
                                    )}
                                    {(field.type === 'text' || field.type === 'number' || field.type === 'password') && (
                                        <Input
                                            className="bg-gray-200/60 border border-gray-200/60"
                                            id={fieldName}
                                            type={field.type}
                                            placeholder={field.placeholder || t("请输入...")}
                                            {...form.register(fieldName, {
                                                required: field.validation?.required && t(`${field.label} 不能为空`),
                                                pattern: field.validation?.pattern ? {
                                                    value: new RegExp(field.validation.pattern),
                                                    message: t("格式不正确")
                                                } : undefined,
                                                onChange: () => form.trigger(fieldName),
                                                onBlur: () => form.trigger(fieldName),
                                            })}
                                        />
                                    )}
                                </>
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

                {/* 修改按钮组的布局 */}
                <div className="lg:flex md:flex hidden justify-start space-x-3">
                    <Button
                        type="submit"
                        variant="surely"
                        className="cursor-pointer"
                        onClick={handleRestart}
                        disabled={loading}
                    >
                        {t('安装')}
                    </Button>

                    <SheetClose className="cursor-pointer border border-gray-200/60 bg-gray-200/60 text-black/70 rounded-md text-sm shadow-sm hover:bg-white hover:border-theme-color/85 hover:text-theme-color/85 h-9 px-5 py-2">
                        {t('取消')}
                    </SheetClose>
                </div>

                {/* 添加小屏幕下的固定按钮组 */}
                <div className="lg:hidden md:hidden flex absolute -top-32 pt-6 right-0  z-50">
                    <Button
                        type="submit"
                        variant="minsure"
                        className="cursor-pointer text-theme-color text-lg font-normal "
                        onClick={handleRestart}
                        disabled={loading}
                    >
                        {t('安装')}
                    </Button>
                </div>
            </form>
        </Form>
    </>
    );
}
