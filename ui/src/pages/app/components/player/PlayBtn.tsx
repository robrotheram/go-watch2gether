export const PlayBtn = ({ status, play, pause }) => {
    if (status === "PLAY") {
        return (
            <button onClick={() => pause()} data-tooltip-target="tooltip-pause" type="button" className="inline-flex items-center justify-center p-2.5 mx-2 font-medium bg-purple-600 rounded-full hover:bg-purple-700 group focus:ring-4 focus:outline-none focus:ring-purple-800">
                <svg className="w-4 h-4 text-white" viewBox="0 0 10 14" fill="currentColor" xmlns="http://www.w3.org/2000/svg" aria-hidden="true">
                    <path fillRule="evenodd" clipRule="evenodd" d="M0.625 1.375C0.625 1.02982 0.904823 0.75 1.25 0.75H2.5C2.84518 0.75 3.125 1.02982 3.125 1.375V12.625C3.125 12.9702 2.84518 13.25 2.5 13.25H1.25C1.08424 13.25 0.925268 13.1842 0.808058 13.0669C0.690848 12.9497 0.625 12.7908 0.625 12.625L0.625 1.375ZM6.875 1.375C6.875 1.02982 7.15482 0.75 7.5 0.75H8.75C8.91576 0.75 9.07473 0.815848 9.19194 0.933058C9.30915 1.05027 9.375 1.20924 9.375 1.375L9.375 12.625C9.375 12.9702 9.09518 13.25 8.75 13.25H7.5C7.15482 13.25 6.875 12.9702 6.875 12.625V1.375Z" fill="currentColor" />
                </svg>
                <span className="sr-only">Pause video</span>
            </button>
        )
    }
    return (
        <button onClick={() => play()} data-tooltip-target="tooltip-pause" type="button" className="inline-flex items-center justify-center p-2.5 mx-2 font-medium bg-purple-600 rounded-full hover:bg-purple-700 group focus:ring-4 focus:outline-none focus:ring-purple-800">
            <svg className="w-4 h-4 text-white" aria-hidden="true" fill="none" stroke="currentColor" strokeWidth="2.5" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                <path d="M5.25 5.653c0-.856.917-1.398 1.667-.986l11.54 6.348a1.125 1.125 0 010 1.971l-11.54 6.347a1.125 1.125 0 01-1.667-.985V5.653z" strokeLinecap="round" strokeLinejoin="round"></path>
            </svg>
            <span className="sr-only">Play video</span>
        </button>
    )
}