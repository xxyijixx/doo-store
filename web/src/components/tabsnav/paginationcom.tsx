import {
    Pagination,
    PaginationContent,
    PaginationItem,
    PaginationLink,
    PaginationNext,
    PaginationPrevious,
} from "@/components/ui/pagination"

interface PaginationComProps {
    currentPage: number;
    totalPages: number;
    onPageChange: (page: number) => void;
}

export function PaginationCom({ currentPage, totalPages, onPageChange }: PaginationComProps) {   
    return (
        
        <Pagination>
            <PaginationContent>
            <PaginationItem>
                <PaginationPrevious 
                    href="#" 
                    disabled={currentPage === 1 } 
                    onClick={() => {
                        if (currentPage > 1) {
                            onPageChange(currentPage - 1);
                        }
                    }} 
                />
            </PaginationItem>
            {[...Array(totalPages)].map((_, index) => {
                    const page = index + 1;
                    return (
                        <PaginationItem key={page}>
                            <PaginationLink 
                                isActive={page === currentPage}
                                onClick={() => onPageChange(page)}
                            >
                                {page}
                            </PaginationLink>
                        </PaginationItem>
                    );
                })}
            <PaginationItem>
                <PaginationNext 
                    href="#" 
                    disabled={currentPage === totalPages} 
                    onClick={() => {
                        if (currentPage < totalPages) {
                            onPageChange(currentPage + 1);
                        }
                    }} 
                />
            </PaginationItem>
            </PaginationContent>
        </Pagination>
        )
}