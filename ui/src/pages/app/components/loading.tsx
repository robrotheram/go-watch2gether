import React from "react";

export const Loading = () => {
    return(
        <div className="h-full w-full flex justify-center items-center ">
            <div className="bg-violet-800 w-48 h-48  absolute animate-ping rounded-full delay-10s shadow-2xl"></div>
            <div className="bg-violet-700 w-32 h-32  absolute animate-ping rounded-full delay-5s shadow-xl"></div>
            <div className="bg-violet-600 w-16 h-16  absolute animate-ping rounded-full shadow-xl"></div>
        </div>
    )
}