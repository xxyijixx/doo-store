
export const LoadingOverlay = () => {
    return (
        <div className=" fixed inset-0 z-50 flex flex-col items-center justify-center bg-black bg-opacity-50">
            <div className="text-white text-2xl">Loading...</div>
            <div className="text-white text-2xl">wait a moment please...</div>
        </div>
    );
};
