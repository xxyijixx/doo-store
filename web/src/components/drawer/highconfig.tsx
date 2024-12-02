
import { Checkbox } from "@/components/ui/checkbox";
import { Textarea } from "@/components/ui/textarea";
import { Label } from "@/components/ui/label";
import { useState } from "react";
import { Input } from "@/components/ui/input";
import { Item } from "@/type.d/common";
import { useTranslation } from "react-i18next";

interface HighConfigProps {
    app: Item;  // 接收 app 数据
    dockerCompose: string;  // dockerCompose 状态
    cpuLimit: string;  // CPU 限制
    memoryLimit: string;  // 内存限制
    setDockerCompose: (value: string) => void;  // 设置 dockerCompose 状态
    setCpuLimit: (value: string) => void;  // 设置 CPU 限制
    setMemoryLimit: (value: string) => void;  // 设置内存限制

}

export function HighConfig(
    { 
        dockerCompose, 
        cpuLimit, 
        memoryLimit, 
        setDockerCompose, 
        setCpuLimit, 
        setMemoryLimit, 
    } : HighConfigProps) {
    
    const { t } = useTranslation();
    const [isAdvancedSettingsEnabled, setAdvancedSettingsEnabled] = useState(false);



    const handleCheckboxChange = () => {
        setAdvancedSettingsEnabled(prev => !prev); // 切换状态
    };


    return (
        <>
            <div className="flex gap-2 text-left">
                <Checkbox 
                    checked={isAdvancedSettingsEnabled} //状态绑定
                    onCheckedChange={handleCheckboxChange}  //处理
                    className="mt-1"
                />
                <Label>{t('高级设置')}</Label>
            </div>
            {isAdvancedSettingsEnabled && ( // 根据状态显示或隐藏文本区域
                <>
                    <Textarea 
                        className="resize-none bg-gray-200/60 border-none" 
                        placeholder={t("输入高级设置...")}
                        value={dockerCompose} // 将返回的数据填充到textarea
                        onChange={(e) => setDockerCompose(e.target.value)} // 如果需要编辑
                        />
                    <div className="w-full">
                                <div className="lg:flex md:flex justify-between py-3 ">
                                    <Label className="mt-2" htmlFor='input'>{t('CPU限制：')}</Label>
                                    <Input 
                                        className='lg:w-5/6 md:w-5/6 bg-gray-200/60 border-gray-200/60'
                                        id="cpuLimit"
                                        value={cpuLimit} // 绑定cpuLimit状态
                                        onChange={(e) => setCpuLimit(e.target.value)} // 更新状态
                                        placeholder="1"
                                        />
                                </div>
                                <div className="lg:flex md:flex justify-between py-3">
                                    <Label className="mt-2"  htmlFor='input'>{t('内存限制：')}</Label>
                                    <Input 
                                        className='lg:w-5/6 md:w-5/6 bg-gray-200/60 border-gray-200/60'
                                        id="memoryLimit"
                                        value={memoryLimit} // 绑定cpuLimit状态
                                        onChange={(e) => setMemoryLimit(e.target.value) } // 更新状态
                                        placeholder="120m 或 12g"
                                        />
                                </div>
                    </div>

                </>
                
            )}
        </>
        
    )
}