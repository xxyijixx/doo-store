import {
    // Pagination,
    // PaginationContent,
    PaginationItem,
    PaginationLink,
    PaginationNext,
    PaginationPrevious,
    // PaginationEllipsis,
} from "@/components/ui/pagination"
import { useState } from "react"

interface PaginationComProps {
    currentPage: number;
    totalPages: number;
    totalItems: number;
    pageSize: number;
    onPageChange: (page: number) => void;
    onPageSizeChange: (pageSize: number) => void;
}

export function PaginationCom({ 
    currentPage, 
    totalPages, 
    totalItems,
    pageSize,
    onPageChange,
}: PaginationComProps) {   
    const [jumpPage, setJumpPage] = useState<string>("");

    const handleJump = () => {
        const page = parseInt(jumpPage);
        if (!isNaN(page) && page > 0 && page <= totalPages) {
            onPageChange(page);
            // setJumpPage("");
        }
    };

    const handleKeyPress = (e: React.KeyboardEvent) => {
        if (e.key === 'Enter') {
            handleJump();
        }
    };

    return (
        <div className="flex items-center justify-end gap-4 text-sm mr-6 lg:mt-6 md:mt-6 mt-6">
            <span>共 {totalItems || 0} 条</span>
            
            <div className="flex items-center list-none">
                <PaginationItem>
                    <PaginationPrevious 
                        href="#" 
                        disabled={currentPage <= 1} 
                        onClick={(e) => {
                            e.preventDefault();
                            if (currentPage > 1) {
                                onPageChange(currentPage - 1);
                            }
                        }} 
                    />
                </PaginationItem>
                
                <PaginationItem>
                    <PaginationLink 
                        href="#"
                        isActive={true}
                    >
                        {currentPage}
                    </PaginationLink>
                </PaginationItem>
                
                <PaginationItem>
                    <PaginationNext 
                        href="#" 
                        disabled={currentPage >= totalPages} 
                        onClick={(e) => {
                            e.preventDefault();
                            if (currentPage < totalPages) {
                                onPageChange(currentPage + 1);
                            }
                        }} 
                    />
                </PaginationItem>
            </div>

            <div className="border border-gray-100 rounded px-2 py-1">
                {pageSize}条/页
            </div>
            
            <div className="items-center gap-2 lg:flex md:flex hidden">
                <span>跳至</span>
                <input
                    type="number"
                    min={1}
                    max={totalPages}
                    value={jumpPage}
                    onChange={(e) => setJumpPage(e.target.value)}
                    onKeyPress={handleKeyPress}
                    onBlur={handleJump}
                    className="border rounded w-12 px-2 mx-0.5 py-1"
                />
                <span>页</span>
            </div>
        </div>
    );
}