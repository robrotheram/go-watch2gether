import { useEffect, useState } from "react";
import Card from "./Cards"
import Header from "./Header"
import Player from "./Player"
import toast, { Toaster, useToaster } from 'react-hot-toast';
import { addVideoController, getController } from "../api/watch2gether";
import { useNavigate } from "react-router-dom";

export const AddVideoCtrl = () => {
    const [video, setVideo] = useState("");
    const addVideo = async () => {
        if (video.length == 0) {
            return
        }

        try {
            await addVideoController(video)
            toast.success("Video is being added to the queue please wait");
        } catch (e) {
            toast.error("Unable to add video: invalid video url");
        }
        setVideo("")
    }
    const handleKeyPress = (e) => {
        if (e.key == 'Enter') {
            addVideo()
        }
    }
    return (

        <div className="w-full bg-violet-950 shadow-2xl z-20 flex items-center">
            <div className="inset-y-0 left-0 flex items-center pl-3 pointer-events-none">
                <svg fill="none" className="w-5 h-5 text-gray-100" stroke="currentColor" strokeWidth="1.5" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg" aria-hidden="true">
                    <path strokeLinecap="round" strokeLinejoin="round" d="M3.375 19.5h17.25m-17.25 0a1.125 1.125 0 01-1.125-1.125M3.375 19.5h1.5C5.496 19.5 6 18.996 6 18.375m-3.75 0V5.625m0 12.75v-1.5c0-.621.504-1.125 1.125-1.125m18.375 2.625V5.625m0 12.75c0 .621-.504 1.125-1.125 1.125m1.125-1.125v-1.5c0-.621-.504-1.125-1.125-1.125m0 3.75h-1.5A1.125 1.125 0 0118 18.375M20.625 4.5H3.375m17.25 0c.621 0 1.125.504 1.125 1.125M20.625 4.5h-1.5C18.504 4.5 18 5.004 18 5.625m3.75 0v1.5c0 .621-.504 1.125-1.125 1.125M3.375 4.5c-.621 0-1.125.504-1.125 1.125M3.375 4.5h1.5C5.496 4.5 6 5.004 6 5.625m-3.75 0v1.5c0 .621.504 1.125 1.125 1.125m0 0h1.5m-1.5 0c-.621 0-1.125.504-1.125 1.125v1.5c0 .621.504 1.125 1.125 1.125m1.5-3.75C5.496 8.25 6 7.746 6 7.125v-1.5M4.875 8.25C5.496 8.25 6 8.754 6 9.375v1.5m0-5.25v5.25m0-5.25C6 5.004 6.504 4.5 7.125 4.5h9.75c.621 0 1.125.504 1.125 1.125m1.125 2.625h1.5m-1.5 0A1.125 1.125 0 0118 7.125v-1.5m1.125 2.625c-.621 0-1.125.504-1.125 1.125v1.5m2.625-2.625c.621 0 1.125.504 1.125 1.125v1.5c0 .621-.504 1.125-1.125 1.125M18 5.625v5.25M7.125 12h9.75m-9.75 0A1.125 1.125 0 016 10.875M7.125 12C6.504 12 6 12.504 6 13.125m0-2.25C6 11.496 5.496 12 4.875 12M18 10.875c0 .621-.504 1.125-1.125 1.125M18 10.875c0 .621.504 1.125 1.125 1.125m-2.25 0c.621 0 1.125.504 1.125 1.125m-12 5.25v-5.25m0 5.25c0 .621.504 1.125 1.125 1.125h9.75c.621 0 1.125-.504 1.125-1.125m-12 0v-1.5c0-.621-.504-1.125-1.125-1.125M18 18.375v-5.25m0 5.25v-1.5c0-.621.504-1.125 1.125-1.125M18 13.125v1.5c0 .621.504 1.125 1.125 1.125M18 13.125c0-.621.504-1.125 1.125-1.125M6 13.125v1.5c0 .621-.504 1.125-1.125 1.125M6 13.125C6 12.504 5.496 12 4.875 12m-1.5 0h1.5m-1.5 0c-.621 0-1.125.504-1.125 1.125v1.5c0 .621.504 1.125 1.125 1.125M19.125 12h1.5m0 0c.621 0 1.125.504 1.125 1.125v1.5c0 .621-.504 1.125-1.125 1.125m-17.25 0h1.5m14.25 0h1.5"></path>
                </svg>
            </div>
            <input onKeyDown={handleKeyPress} type="text" value={video} onChange={(e) => setVideo(e.target.value)} className="block w-full p-4 pl-10 text-xl bg-transparent text-white focus:ring-purple-500 focus:border-purple-500" placeholder="Add New video" required />
            <button onClick={() => addVideo()} className="text-white whitespace-nowrap py-2 px-4 mr-2 bg-purple-700 hover:bg-purple-800 focus:ring-4 focus:outline-none focus:ring-purple-300 font-medium rounded-lg">Add Video</button>
        </div>

    )
}

const Notifications = () => {
    const { toasts } = useToaster();
    

    const SusccessIcon = (
        <div class="inline-flex items-center justify-center flex-shrink-0 w-8 h-8 text-green-500 bg-green-100 rounded-lg dark:bg-green-800 dark:text-green-200">
            <svg aria-hidden="true" class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20" xmlns="http://www.w3.org/2000/svg"><path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd"></path></svg>
            <span class="sr-only">Check icon</span>
        </div>
    )

    const WarningIcon = (
        <div class="inline-flex items-center justify-center flex-shrink-0 w-8 h-8 text-orange-500 bg-orange-100 rounded-lg dark:bg-orange-700 dark:text-orange-200">
            <svg aria-hidden="true" class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20" xmlns="http://www.w3.org/2000/svg"><path fill-rule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clip-rule="evenodd"></path></svg>
            <span class="sr-only">Warning icon</span>
        </div>
    )



    return (
        <div className="absolute z-50 top-2 right-2 left-2 sm:left-auto md:left-auto">
            {toasts.filter((toast) => toast.visible)
                .map((t) => {
                    setTimeout(() => toast.dismiss(t.id), t.duration)
                    return (
                        <div key={t.id} style={{ opacity: t.visible ? 1 : 0 }}  {...t.ariaProps} id="toast-default" class="flex items-center w-full p-4 mb-2 rounded-xl text-white bg-zinc-800" role="alert">
                            {t.type === "success" ? SusccessIcon : WarningIcon}
                            <div class="ml-3 text-md font-normal px-5">{t.message}</div>
                            <button onClick={() => toast.dismiss(t.id)} type="button" class="ml-auto -mx-1.5 -my-1.5 hover:text-purple-900 rounded-lg focus:ring-2 focus:ring-gray-300 p-1.5 hover:bg-gray-100 inline-flex h-8 w-8 dark:text-white dark:hover:text-white dark:bg-purple-800 dark:hover:bg-purple-700" data-dismiss-target="#toast-default" aria-label="Close">
                                <span class="sr-only">Close</span>
                                <svg aria-hidden="true" class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20" xmlns="http://www.w3.org/2000/svg"><path fill-rule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clip-rule="evenodd"></path></svg>
                            </button>
                        </div>
                    )
                }
                )
            }
        </div>
    );

};

const Controller = () => {
    const navigate = useNavigate();
    const [loading, setLoading] = useState(true)
    const [error, setError] = useState("")
    const [state, setState] = useState({
        Current: {
            duration: 0
        }, Queue: []
    })

    useEffect(() => {
        const interval = setInterval(async () => {
            try {
                let _state = await getController()
                setState(_state)
            } catch (e) {
                navigate(`/app?error`)
            }
            setLoading(false)

        }, 1000);
        return () => clearInterval(interval);
    }, []);


    if (loading) {
        return (
            <div className="min-h-screen w-full flex justify-center items-center bg-black">
                <div className="bg-violet-800 w-48 h-48  absolute animate-ping rounded-full delay-10s shadow-2xl"></div>
                <div className="bg-violet-700 w-32 h-32  absolute animate-ping rounded-full delay-5s shadow-xl"></div>
                <div className="bg-violet-600 w-16 h-16  absolute animate-ping rounded-full shadow-xl"></div>
            </div>
        )
    }

    return (
        <div className="flex flex-col w-full h-full">
            <Notifications />
            <AddVideoCtrl />
            <div className='bg-violet-800 w-full h-full flex flex-col' style={{ "overflow": "auto" }}>
                {state.Current.id && <Header current={state.Current} />}
                <div className='w-full flex-grow' >
                    <div className='w-full h-full shadow-body px-4 md:px-10 text-white'>
                        <Card queue={state.Queue} />
                    </div>
                </div>
                <Player state={state} />
            </div>
        </div>
    )
}
export default Controller