import { loopVideoController, pauseVideoController, playVideoController, shuffleVideoController, skipVideoController } from "../api/watch2gether";



const Player = ({ state }) => {
    const formatTime = (seconds) => {
        if (seconds === undefined) {
            seconds = 0
        }
        let iso = new Date(seconds / 1000000).toISOString()
        return iso.substring(11, iso.length - 5);
    }
    const playerProgress = (current, total) => {
        let pct = current / total * 100
        return Math.min(Math.max(pct, 0), 100)
    }

    const handleShuffle = () => {
        shuffleVideoController();
    }
    const handlePlay = () => {
        playVideoController();
    }
    const handlePause = () => {
        pauseVideoController();
    }
    const handleSkip = () => {
        skipVideoController();
    }
    const handleLoop = () => {
        loopVideoController();
    }




    return (
        <div className="fixed bottom-0 left-0 z-50 grid w-full h-24 grid-cols-1 px-8  md:grid-cols-3 bg-zinc-900">
            <div></div>
            {state.Active ?
                <div className="flex items-center w-full">
                    <div className="w-full">
                        <div className="flex items-center justify-center mx-auto mb-1">
                            <div className="p-2.5 w-5" />
                            <button onClick={() => handleShuffle()} data-tooltip-target="tooltip-shuffle" type="button" className="p-2.5 group rounded-full hover:bg-gray-100 mr-1 focus:outline-none focus:ring-4 focus:ring-gray-200 dark:focus:ring-gray-600 dark:hover:bg-gray-600">
                                <svg className="w-5 h-5 text-gray-500 dark:text-gray-300 group-hover:text-gray-900 dark:group-hover:text-white" fill="currentColor" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 512 512" aria-hidden="true"><path d="M403.8 34.4c12-5 25.7-2.2 34.9 6.9l64 64c6 6 9.4 14.1 9.4 22.6s-3.4 16.6-9.4 22.6l-64 64c-9.2 9.2-22.9 11.9-34.9 6.9s-19.8-16.6-19.8-29.6V160H352c-10.1 0-19.6 4.7-25.6 12.8L284 229.3 244 176l31.2-41.6C293.3 110.2 321.8 96 352 96h32V64c0-12.9 7.8-24.6 19.8-29.6zM164 282.7L204 336l-31.2 41.6C154.7 401.8 126.2 416 96 416H32c-17.7 0-32-14.3-32-32s14.3-32 32-32H96c10.1 0 19.6-4.7 25.6-12.8L164 282.7zm274.6 188c-9.2 9.2-22.9 11.9-34.9 6.9s-19.8-16.6-19.8-29.6V416H352c-30.2 0-58.7-14.2-76.8-38.4L121.6 172.8c-6-8.1-15.5-12.8-25.6-12.8H32c-17.7 0-32-14.3-32-32s14.3-32 32-32H96c30.2 0 58.7 14.2 76.8 38.4L326.4 339.2c6 8.1 15.5 12.8 25.6 12.8h32V320c0-12.9 7.8-24.6 19.8-29.6s25.7-2.2 34.9 6.9l64 64c6 6 9.4 14.1 9.4 22.6s-3.4 16.6-9.4 22.6l-64 64z" fill="currentColor" /></svg>
                                <span className="sr-only">Shuffle video</span>
                            </button>
                            <div id="tooltip-shuffle" role="tooltip" className="absolute z-10 invisible inline-block px-3 py-2 text-sm font-medium text-white transition-opacity duration-300 bg-gray-900 rounded-lg shadow-sm opacity-0 tooltip dark:bg-gray-700">
                                Shuffle video
                                <div className="tooltip-arrow" data-popper-arrow></div>
                            </div>

                            {state.State !== "PAUSED" && state.State !=="STOPPED"  ?
                                <button onClick={() => handlePause()} data-tooltip-target="tooltip-pause" type="button" className="inline-flex items-center justify-center p-2.5 mx-2 font-medium bg-purple-600 rounded-full hover:bg-purple-700 group focus:ring-4 focus:ring-purple-300 focus:outline-none dark:focus:ring-purple-800">
                                    <svg className="w-4 h-4 text-white" viewBox="0 0 10 14" fill="currentColor" xmlns="http://www.w3.org/2000/svg" aria-hidden="true">
                                        <path fillRule="evenodd" clipRule="evenodd" d="M0.625 1.375C0.625 1.02982 0.904823 0.75 1.25 0.75H2.5C2.84518 0.75 3.125 1.02982 3.125 1.375V12.625C3.125 12.9702 2.84518 13.25 2.5 13.25H1.25C1.08424 13.25 0.925268 13.1842 0.808058 13.0669C0.690848 12.9497 0.625 12.7908 0.625 12.625L0.625 1.375ZM6.875 1.375C6.875 1.02982 7.15482 0.75 7.5 0.75H8.75C8.91576 0.75 9.07473 0.815848 9.19194 0.933058C9.30915 1.05027 9.375 1.20924 9.375 1.375L9.375 12.625C9.375 12.9702 9.09518 13.25 8.75 13.25H7.5C7.15482 13.25 6.875 12.9702 6.875 12.625V1.375Z" fill="currentColor" />
                                    </svg>
                                    <span className="sr-only">Pause video</span>

                                </button>
                                : <button onClick={() => handlePlay()} data-tooltip-target="tooltip-pause" type="button" className="inline-flex items-center justify-center p-2.5 mx-2 font-medium bg-purple-600 rounded-full hover:bg-purple-700 group focus:ring-4 focus:ring-purple-300 focus:outline-none dark:focus:ring-purple-800">
                                    <svg className="w-4 h-4 text-white" aria-hidden="true" fill="none" stroke="currentColor" strokeWidth="2.5" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                                        <path d="M5.25 5.653c0-.856.917-1.398 1.667-.986l11.54 6.348a1.125 1.125 0 010 1.971l-11.54 6.347a1.125 1.125 0 01-1.667-.985V5.653z" strokeLinecap="round" strokeLinejoin="round"></path>
                                    </svg>
                                    <span className="sr-only">Play video</span>
                                </button>
                            }
                            <button onClick={() => handleSkip()} data-tooltip-target="tooltip-next" type="button" className="p-2.5 group rounded-full hover:bg-gray-100 mr-1 focus:outline-none focus:ring-4 focus:ring-gray-200 dark:focus:ring-gray-600 dark:hover:bg-gray-600">
                                <svg className="w-5 h-5 text-gray-500 dark:text-gray-300 group-hover:text-gray-900 dark:group-hover:text-white" fill="currentColor" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 320 512" aria-hidden="true"><path d="M52.5 440.6c-9.5 7.9-22.8 9.7-34.1 4.4S0 428.4 0 416V96C0 83.6 7.2 72.3 18.4 67s24.5-3.6 34.1 4.4l192 160L256 241V96c0-17.7 14.3-32 32-32s32 14.3 32 32V416c0 17.7-14.3 32-32 32s-32-14.3-32-32V271l-11.5 9.6-192 160z" fill="currentColor" /></svg>
                                <span className="sr-only">Next video</span>
                            </button>
                            <div id="tooltip-next" role="tooltip" className="absolute z-10 invisible inline-block px-3 py-2 text-sm font-medium text-white transition-opacity duration-300 bg-gray-900 rounded-lg shadow-sm opacity-0 tooltip dark:bg-gray-700">
                                Next video
                                <div className="tooltip-arrow" data-popper-arrow></div>
                            </div>
                            <button onClick={() => handleLoop()} data-tooltip-target="tooltip-restart" type="button" className="relative p-2.5 group rounded-full hover:bg-gray-100 mr-1 focus:outline-none focus:ring-4 focus:ring-gray-200 dark:focus:ring-gray-600 dark:hover:bg-gray-600">
                                <svg className="w-5 h-5 text-gray-500 dark:text-gray-300 group-hover:text-gray-900 dark:group-hover:text-white" fill="currentColor" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 512 512" aria-hidden="true"><path d="M0 224c0 17.7 14.3 32 32 32s32-14.3 32-32c0-53 43-96 96-96H320v32c0 12.9 7.8 24.6 19.8 29.6s25.7 2.2 34.9-6.9l64-64c12.5-12.5 12.5-32.8 0-45.3l-64-64c-9.2-9.2-22.9-11.9-34.9-6.9S320 19.1 320 32V64H160C71.6 64 0 135.6 0 224zm512 64c0-17.7-14.3-32-32-32s-32 14.3-32 32c0 53-43 96-96 96H192V352c0-12.9-7.8-24.6-19.8-29.6s-25.7-2.2-34.9 6.9l-64 64c-12.5 12.5-12.5 32.8 0 45.3l64 64c9.2 9.2 22.9 11.9 34.9 6.9s19.8-16.6 19.8-29.6V448H352c88.4 0 160-71.6 160-160z" fill="currentColor" /></svg>
                                <span className="sr-only">Loop video</span>
                                {state.Loop && <div className="absolute bottom-1 right-1 w-3 h-3 text-xs font-bold text-white bg-red-500 border-2 border-white rounded-full  dark:border-gray-900" />}
                            </button>
                        </div>
                        <div className="flex items-center justify-between space-x-2">
                            <span className="text-sm font-medium text-gray-500 dark:text-gray-400">{formatTime(state.Proccessing)}</span>
                            <div className="w-full bg-gray-200 rounded-full h-1.5 dark:bg-gray-800">
                                <div className="bg-purple-600 h-1.5 rounded-full" style={{ "width": `${playerProgress(state.Proccessing, state.Current.duration)}%` }}></div>
                            </div>
                            <span className="text-sm font-medium text-gray-500 dark:text-gray-400">{formatTime(state.Current.duration)}</span>
                        </div>
                    </div>
                </div>
            :
            <div className="flex items-center w-full text-center text-white">
                <div className="w-full text-center">Player is not active. <br/> Join the bot to one of the voice channels</div>
            </div>
            }
        </div>
    )
}

export default Player