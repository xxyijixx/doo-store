import { Sheet, SheetContent, SheetDescription, SheetHeader, SheetTitle } from "@/components/ui/sheet";
import { EditForm } from "@/components/drawer/editdraform";
import { Item } from "@/type.d/common";
import { useTranslation } from "react-i18next";
import { ChevronLeftIcon } from "@radix-ui/react-icons"

interface EditDrawerProps {
    isOpen: boolean;
    onClose: () => void;
    app: Item;
}

function EditDrawer ({ isOpen, onClose, app }: EditDrawerProps){


    const editSuccess = () => {
        console.log('editSuccess');
    }
    
    const editFalse = () => {
        console.log('editFalse');
    }

    const { t } = useTranslation();

    return (
        <Sheet open={isOpen} onOpenChange={onClose}>
            <SheetContent className='lg:overflow-y-auto md:overflow-auto overflow-auto'>
                <SheetHeader>
                    <SheetTitle className='lg:ml-0 md:ml-0 pl-2  text-gray-700 z-50 lg:bg-transparent md:bg-transparent bg-gray-200/50 lg:py-0 md:py-0 py-3 flex items-center gap-2'>
                        <ChevronLeftIcon 
                            className="h-6 w-6 lg:hidden md:hidden block"
                            onClick={() => onClose()}
                        />
                        {t('参数修改')}
                    </SheetTitle>
                    <hr  className='lg:block md:block hidden'/>
                </SheetHeader>
                <SheetDescription className='pt-3'>
                </SheetDescription>
                {/* 将 app 和回调传递给 EditForm */}
                <EditForm 
                    app={app} 
                    onEditSuccess={editSuccess} 
                    onEditFalse={editFalse} 
                />
            </SheetContent>
        </Sheet>
    );
};

export default EditDrawer;
