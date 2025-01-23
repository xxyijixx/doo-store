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
import { ChevronLeftIcon } from "@radix-ui/react-icons"


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
            return 'border border-theme-color text-sm font-normal text-theme-color hover:text-theme-color/80 hover:border-theme-color/80 h-8 px-3 whitespace-nowrap cursor-pointer'; // 未使用状态
        }
        return 'border-2 border-theme-color text-theme-color'; // 默认样式
    };


    //handleInstallSuccess执行安装操作后
    const handleInstallSuccess = () => {
        setButtonText(t("安装"));
        setCurrentStatus('InUse');
        setIsOpen(false);
        
        // 直接切换到已安装页面
        window.dispatchEvent(new CustomEvent('switchToInstalled'));
    };

    //handleInstallFalse失败运行安装局部更新状态
    const handleInstallFalse = () => {
        toast({
            title: t("安装失败"),
            description: t("安装过程中发生错误，请重试。"),
            variant: "destructive",
            duration: 3000,
            className: "fixed top-16  lg:top-3 md:top-3 lg:right-6  md:right-4 right-1/2 translate-x-1/2 lg:translate-x-0 md:translate-x-0 w-[350px]"
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
                                e.preventDefault();
                                e.stopPropagation();
                            } else {
                                handleInstallClick(); // 处理安装逻辑
                            }
                        }}
                        className={` rounded-md pt-1.5 ${getButtonStyles()}`}>
                        {buttonText}
                    </div>
                </SheetTrigger>
                
                <SheetContent className='lg:overflow-y-auto md:overflow-auto overflow-auto'>
                    <SheetHeader>
                        <SheetTitle className='lg:ml-2 md:ml-2  text-gray-700 z-50 lg:bg-transparent md:bg-transparent bg-gray-200/50 lg:py-0 md:py-0 py-3 flex items-center gap-2'>
                            <ChevronLeftIcon 
                                className="h-6 w-6 ml-3 lg:hidden md:hidden block"
                                onClick={() => setIsOpen(false)}
                                />
                            {t('返回')}
                        </SheetTitle>
                        <hr  className='lg:block md:block hidden'/>
                        <SheetDescription>
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

