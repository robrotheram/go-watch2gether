import { useState, memo } from "react"
import { useSwipeable } from "react-swipeable";
import media_svg from "../../../assets/media.svg"
import './card.css'

const LoadingSkeleton = () => (
    <div className="flex-shrink-0 transform transition-transform">
        <div className="h-20 w-20 md:h-32 md:w-48 rounded-lg bg-violet-800/50 relative overflow-hidden">
            <div className="absolute inset-0 bg-gradient-to-r from-transparent via-violet-600/20 to-transparent skeleton-loading" />
        </div>
    </div>
);

const QueueItem = memo(({ key, pos, video, children }) => {
    const [showMenu, setShowMenu] = useState(false)
    const handlers = useSwipeable({
        onSwipedLeft: () => setShowMenu(true),
        onSwipedRight: () => setShowMenu(false),
    });

    const formatTime = (seconds) => {
        if (seconds === undefined) {
            seconds = 0
        }
        let iso = new Date(seconds / 1000000).toISOString()
        return iso.substring(11, iso.length - 5);
    }

    return (
        <li key={key} className="py-3" {...handlers} onMouseEnter={() => setShowMenu(true)} onMouseLeave={() => setShowMenu(false)}>
            <div className="flex items-center space-x-4 md:space-x-10">
                <div className="flex text-md text-white text-center" style={{ width: "20px" }}>
                    {pos}
                </div>
                <div className="flex-shrink-0">
                    {video.loading ? (
                        <LoadingSkeleton />
                    ) : (
                        video.thumbnail ?
                            <img loading="lazy" className="h-20 w-20 md:h-32 md:w-48 object-cover rounded-lg" src={video.thumbnail} alt={video.title} />
                            :
                            <img loading="lazy" className="h-20 w-20 md:h-32 md:w-48 rounded-lg p-4" src={media_svg} alt={video.title} />
                    )}
                </div>
                <div className="flex-1 min-w-0 md:flex md:justify-aroud">
                    <p className="text-md font-medium truncate text-white w-full md:w-1/3 text-left md:text-center">
                        {video.title}
                    </p>
                    <p className="text-md truncate text-white w-full md:w-1/3 text-left md:text-center">
                        Added by: {video.user}
                    </p>
                    <p className="text-md truncate text-white w-full md:w-1/3 text-left md:text-center">
                        {!video.loading && formatTime(video.time.duration)}
                    </p>
                </div>

                <div className="inline-flex items-center text-base font-semibold text-white">
                    {!video.loading && showMenu && children}
                </div>
            </div>
        </li>
    )
});

const Card = ({ queue, updateQueue }) => {
    const deleteVideo = (item) => {
        const videoList = [...queue];
        const i = videoList.indexOf(item);
        videoList.splice(i, 1);
        updateQueue(videoList);
    };

    const voteUp = (item) => {
        const videoList = [...queue];
        const i = videoList.indexOf(item);
        const z = videoList[i - 1];
        videoList[i - 1] = videoList[i];
        videoList[i] = z;
        updateQueue(videoList);
    };

    const moveToTop = (item) => {
        let videoList = [...queue];
        const i = videoList.indexOf(item);
        videoList.splice(i, 1);
        videoList = [item, ...videoList];
        updateQueue(videoList);
    };

    const voteDown = (item) => {
        const videoList = [...queue];
        const i = videoList.indexOf(item);
        const z = videoList[i + 1];
        videoList[i + 1] = videoList[i];
        videoList[i] = z;
        updateQueue(videoList);
    };


    if (queue.length > 0) {
        return (
            <ul role="list" className="divide-y divide-gray-200 mt-8">
                {queue.map((video, i) => <QueueItem key={video.qid} pos={i + 1} video={video}>
                    <div className="inline-flex rounded-md shadow-sm" role="group">
                        <button type="button" onClick={() => moveToTop(video)} className="inline-flex items-center px-4 py-2 text-sm font-medium text-gray-900 border border-gray-200 rounded-l-lg hover:bg-gray-100 hover:text-purple-700 focus:z-10 focus:ring-2 focus:ring-purple-700 focus:text-purple-700 bg-zinc-100 ">
                            <svg aria-hidden="true" className=" h-6 fill-current" fill="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                                <path strokeLinecap="round" strokeLinejoin="round" d="M4.5 12.75l7.5-7.5 7.5 7.5m-15 6l7.5-7.5 7.5 7.5"></path>
                            </svg>
                        </button>
                        <button type="button" onClick={() => voteUp(video)} className="inline-flex items-center px-4 py-2 text-sm font-medium text-gray-900 bg-white border-t border-b border-gray-200 hover:bg-gray-100 hover:text-purple-700 focus:z-10 focus:ring-2 focus:ring-purple-700 focus:text-purple-700 ">
                            <svg className="h-6" fill="none" stroke="currentColor" strokeWidth="1.5" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg" aria-hidden="true">
                                <path strokeLinecap="round" strokeLinejoin="round" d="M4.5 15.75l7.5-7.5 7.5 7.5"></path>
                            </svg>
                        </button>
                        <button type="button" onClick={() => voteDown(video)} className="inline-flex items-center px-4 py-2 text-sm font-medium text-gray-900 bg-white border-t border-b border-gray-200 hover:bg-gray-100 hover:text-purple-700 focus:z-10 focus:ring-2 focus:ring-purple-700 focus:text-purple-700">
                            <svg className="h-6" fill="none" stroke="currentColor" strokeWidth="1.5" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg" aria-hidden="true">
                                <path strokeLinecap="round" strokeLinejoin="round" d="M19.5 8.25l-7.5 7.5-7.5-7.5"></path>
                            </svg>
                        </button>
                        <button type="button" onClick={() => deleteVideo(video)} className="inline-flex items-center px-4 py-2 text-sm font-medium text-gray-900 bg-white border border-gray-200 rounded-r-md hover:bg-gray-100 hover:text-purple-700 focus:z-10 focus:ring-2 focus:ring-purple-700 focus:text-purple-700">
                            <svg fill="none" className="h-5" stroke="currentColor" strokeWidth="1.5" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg" aria-hidden="true">
                                <path strokeLinecap="round" strokeLinejoin="round" d="M14.74 9l-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 01-2.244 2.077H8.084a2.25 2.25 0 01-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 00-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 013.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 00-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 00-7.5 0"></path>
                            </svg>
                        </button>
                    </div>
                </QueueItem>)}
            </ul>
        )
    }


    return (
        <div className="flex flex-col items-center justify-center py-12 backdrop-blur-sm">
            <div className="p-16 rounded-lg bg-violet-800/20 text-center border border-violet-300 shadow-lg">
                <svg
                    className="mx-auto h-24 w-24 text-violet-400 mb-4"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                    xmlns="http://www.w3.org/2000/svg"
                >
                    <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth={1.5}
                        d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z"
                    />
                    <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth={1.5}
                        d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                    />
                </svg>
                <h3 className="text-2xl font-semibold text-white mb-2">Queue is Empty</h3>
                <p className="text-violet-200">Add some videos to get started!</p>
            </div>
        </div>
    )

}

export default Card