import {
    Sheet,
    SheetContent,
    SheetHeader,
    SheetTitle,
    SheetDescription,
    SheetTrigger,
} from "@/components/ui/sheet"
import { ProfileForm } from "@/components/drawer/draform"
import { useState } from "react";
import { Item } from "@/type.d/common";
import { useToast } from "@/hooks/use-toast";
import { FalseToaster, SuccessToaster } from '../ui/toaster';
import { useTranslation } from "react-i18next";

interface DrawerProps {
    status: string;
    isOpen: boolean;
    app: Item;
    loadData: () => void;
}

function Drawer({ status, app }: DrawerProps) {

    const { t } = useTranslation();
    const [isOpen, setIsOpen] = useState(false);
    const [buttonText, setButtonText] = useState(status === 'InUse' ? t('安装') : t('安装')); // 用状态管理按钮文字
    const [currentStatus, setCurrentStatus] = useState(status); // 用状态管理当前的安装状态，便于控制按钮样式
    const { toast } = useToast();

    // 当前 variant 状态，用于控制渲染
    const [variantState, setVariantState] = useState<"success" | "destructive" | null>(null);

    const handleInstallClick = () => {
        if (status === 'Unused') {
            setIsOpen(true); // 打开侧边栏
        }

    };

    const getButtonStyles = () => {
        if (currentStatus === 'InUse') {
            return 'border-2 border-gray-300 rounded-md bg-gray-300 text-sm text-white shadow-sm h-8 px-3 whitespace-nowrap cursor-not-allowed'; // 在使用状态
        } else if (currentStatus === 'Unused') {
            return 'border-2 border-theme-color text-sm text-theme-color hover:text-theme-color/80 hover:border-theme-color/80 h-8 px-3 whitespace-nowrap cursor-pointer'; // 未使用状态
        }
        return 'border-2 border-theme-color text-theme-color'; // 默认样式
    };


    //handleInstallSuccess成功运行安装局部更新状态
    const handleInstallSuccess = () => {
        toast({
            title: t("安装成功"),
            description: t("应用已成功安装"),
            variant: "success",
            duration: 3000,
        });
        setVariantState("success"); 
        setButtonText(t("安装"));
        setCurrentStatus('InUse'); // 更新状态为 "InUse"（已安装），以改变按钮样式
        setIsOpen(false); // 关闭侧边栏

    };

    //handleInstallFalse失败运行安装局部更新状态
    const handleInstallFalse = () => {
        toast({
            title: t("安装失败"),
            description: t("安装过程中发生错误，请重试。"),
            variant: "destructive",
            duration: 3000,
        });
        setVariantState("destructive");
        setIsOpen(false); // 关闭侧边栏
    };



    return (
        <>
            {variantState === "success" && <SuccessToaster />}
            {variantState === "destructive" && <FalseToaster />}
            <Sheet open={isOpen} onOpenChange={setIsOpen}>
                <SheetTrigger asChild>
                    <div
                        onClick={(e) => {
                            if (status === 'InUse') {
                                e.preventDefault(); // 阻止默认行为
                                e.stopPropagation(); // 阻止事件冒泡
                            } else {
                                handleInstallClick(); // 处理安装逻辑
                            }
                        }}
                        className={` rounded-md pt-1 ${getButtonStyles()}`}>
                        {buttonText}
                    </div>
                </SheetTrigger>
                
                <SheetContent className='lg:overflow-y-auto md:overflow-auto overflow-auto '>
                    <SheetHeader>
                        <SheetTitle className='ml-9 -mt-1.5 text-gray-700'>{t('返回')}</SheetTitle>
                        <hr />
                        <SheetDescription className='pt-3'>
                        </SheetDescription>
                        <ProfileForm
                            app={app}
                            onInstallSuccess={handleInstallSuccess}
                            onFalse={handleInstallFalse}
                        />
                    </SheetHeader>
                </SheetContent>
            </Sheet>
        </>
    )
}



export default Drawer

