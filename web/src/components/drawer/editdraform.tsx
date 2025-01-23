/* eslint-disable @typescript-eslint/no-explicit-any */
/* eslint-disable @typescript-eslint/no-unused-vars */
"use client"

import { useForm } from "react-hook-form";
import { Button } from "@/components/ui/button";
import { Form, FormControl, FormItem, FormLabel, FormMessage } from "@/components/ui/form";
import { SheetClose } from "@/components/ui/sheet";
import { Input } from "@/components/ui/input";
import { Checkbox } from "@/components/ui/checkbox";
import { Label } from "@/components/ui/label";
import { EditHighConfig } from '@/components/drawer/edithighconfig';
import { Item} from "@/type.d/common";
import { useEffect, useState } from "react";
import * as http from "@/api/modules/fouceinter"
import { useTranslation } from "react-i18next";
import { FormField } from "@/api/interface/common"
import { useToast } from "@/hooks/use-toast";
import { SuccessToaster, FalseToaster } from "@/components/ui/toaster"

interface EditProps {
    app: Item;  // 接收 app 数据
    onEditSuccess: () => void; // 新增回调
    onEditFalse: () => void; // 失败回调
}

type FormValues = {
    [key: string]: string | number | boolean; // 根据实际情况调整类型
};

export function EditForm({ app, onEditSuccess, onEditFalse }: EditProps) {
    const { t } = useTranslation();
    const [dockerCompose, setDockerCompose] = useState<string>("");  // 存储docker_compose内容
    const [cpuLimit, setCpuLimit] = useState<string>("1");  // 默认值为 1
    const [memoryLimit, setMemoryLimit] = useState<string>("0");  // 默认值为 120M
    const [loading, setLoading] = useState<boolean>(false);  // 加载状态
    const [error, setError] = useState<string>("");  // 错误信息
    const [formFields, setFormFields] = useState<FormField[]>([]);  // 存储 form_fields 数据
    // 当前 variant 状态，用于控制渲染
    const [variantState, setVariantState] = useState<"success" | "destructive" | null>(null);

    const { toast } = useToast();


    const form = useForm<FormValues>({
        defaultValues: {},
    });

    const {
        setValue,
        handleSubmit,
        formState: { errors }
    } = form;

    // 发起请求获取 form_fields 内容
    useEffect(() => {
        if (app.id) {
            setLoading(true); // 开始加载
            setError(""); // 清空之前的错误
            http.getInsParams(app.id)
                .then((response) => {
                    // if (!response.state===200) {
                    //     throw new Error("请求失败");
                    // }
                    return response;
                })
                .then((data) => {
                    // console.log(data.data);
                    setDockerCompose(data.data?.docker_compose || "");
                    setFormFields(data.data?.params || []);
                    setCpuLimit(data.data?.cpus || cpuLimit);  // 确保这里更新了最新的值
                    setMemoryLimit(data.data?.memory_limit || memoryLimit);  // 同上
                    data.data?.params.forEach(
                        (field) => {
                            const fieldName = field.env_key;
                            setValue(fieldName, field.default ||""); // 设置每个字段的默认值
                        }
                    );
                })
                .catch((error ) => {
                    console.log(error);
                    setError(t("请求失败，请稍后重试"));
                    
                });
                setLoading(false);
        }
    }, [app.id, setValue]);

    const handleRestart = async () => {
        setLoading(true);
        setError(""); // 清除之前的错误信息

        // 校验表单，确保字段不为空
        const isValid = await form.trigger(); // 校验所有表单字段
        if (!isValid) {
            console.log(error)
            setError(t("请填完整的字段信息！"));
            setLoading(false);
            return;
        }

        const formData = form.getValues();
        const params: { [key: string]: string | number | boolean } = {};
        formFields.forEach((field) => {
            params[field.env_key] = formData[field.env_key]
        });

        try {

            // 构建请求体
            const requestBody = {
                cpus: cpuLimit,
                docker_compose: dockerCompose,
                memory_limit: memoryLimit,
                params: params,
            };

            const response = await http.putInsParams(app.id, requestBody)

            const result = await response;
            console.log("请展出",result);
            if (response) {
                // 成功后，更新状态，使得页面渲染新的内容
                console.log("真棒！修改成功");
                toast({
                    title: t("修改成功"),
                    description: t("内容已做修改~"),
                    variant: "success",
                    duration: 2000,
                });
                setVariantState("success");
                onEditSuccess();
                setCpuLimit(result.data?.cpus || cpuLimit);  // 确保获取到最新的值
                setMemoryLimit(result.data?.memory_limit || memoryLimit);  // 同上

                // 更新动态字段
                setFormFields(result.data?.params || []);
                result.data?.params.forEach((field) => {
                    const fieldName = field.env_key;
                    setValue(fieldName, field.default || ""); // 设置每个字段的默认值
                });
            } else {
                console.error("API 请求失败:", result);
                toast({
                    title: t("修改失败"),
                    description: t("内容修改失败，请重试~"),
                    variant: "destructive",
                    duration: 2000,
                });
                setVariantState("destructive");
                
            }
        } catch (error) {
            setError(t("网络错误，请检查网络连接"));
            onEditFalse();
            toast({
                title: t("修改失败"),
                description: t("内容修改失败，请重试~"),
                variant: "destructive",
                duration: 2000,
            });
            setVariantState("destructive");
        } finally {
            setLoading(false);
        }
    };

    return (
        <>
        {variantState === "success" && <SuccessToaster />}
        {variantState === "destructive" && <FalseToaster />}
        <Form {...form} >
            <form className="space-y-8 relative overflow-visible lg:px-0 md:px-0 px-3 pb-3" onSubmit={handleSubmit(handleRestart)}>
                {/* 动态渲染 form_fields */}
                {formFields.sort((a, b) => a.order - b.order).map((field, index) => {
                    const fieldName = field.env_key;
                    
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
                    <EditHighConfig
                        app={app}
                        dockerCompose={dockerCompose}
                        cpuLimit={cpuLimit}  // 传递 cpuLimit
                        memoryLimit={memoryLimit}  // 传递 memoryLimit
                        setDockerCompose={setDockerCompose}
                        setCpuLimit={setCpuLimit}
                        setMemoryLimit={setMemoryLimit}
                    />
                </FormItem>

                <div className="lg:flex md:flex hidden justify-start space-x-3">
                    <Button
                        type="submit"
                        variant="surely"
                        onClick={handleRestart}
                        disabled={loading}
                        >{t('重启')}</Button>
                    <SheetClose
                        className="cursor-pointer border border-gray-200/60 rounded-md bg-gray-200/60 text-sm text-gray-600 shadow-sm hover:bg-white  hover:border-theme-color/85 hover:text-theme-color/85 h-9 px-5 py-2"
                    >{t('取消')}</SheetClose>
                </div>

                 {/* 添加小屏幕下的固定按钮组 */}
                <div className="lg:hidden md:hidden flex absolute -top-28 pt-6 right-0  z-50">
                    <Button
                        type="submit"
                        variant="minsure"
                        className="cursor-pointer text-theme-color text-lg font-normal "
                        onClick={handleRestart}
                        disabled={loading}
                    >
                        {t('重启')}
                    </Button>
                </div>
            </form>
        </Form>
    </>
    );
}
