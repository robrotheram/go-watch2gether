import { useEffect, useRef, useState } from "react"
import { toast } from "react-hot-toast"
import Card from "./components/card"
import { createPlaylist, deletePlaylist, getChannelPlaylists, loadFromPlaylist, updatePlaylist } from "./watch2gether"
import { useOnClickOutside } from "./components/nav"
import { Link } from "react-router-dom"
import { AddVideoCtrl } from "./components/header/Controls"


const ManagePlaylist = ({ playlist, onUpdate }) => {
  const addVideo = async (video) => {
    playlist.videos.push({ url: video })
    try {
      await updatePlaylist(playlist)
      onUpdate()
      toast.success("Video is being added to the playlist please wait");
    } catch (e) {
      toast.error("Unable to add video: invalid video url");
    }
  }
  const updateQueue = async (queue) => {
    playlist.videos = queue
    try {
      await updatePlaylist(playlist)
      onUpdate()
      toast.success("Playlist updated")
    } catch (e) {
      console.error(e)
      toast.error("Sorry there was an issue updating the playlist")
    }
  }
  return (
    <div className="w-full absolute top-64 bottom-0 md:relative md:top-0 overflow-auto">
      <div className="sticky w-full top-0"><AddVideoCtrl onAddVideo={addVideo} /></div>
      <div className="px-8" >
        <Card queue={playlist.videos} updateQueue={updateQueue} />
      </div>
    </div>
  )
}

const PlaylistHeader = ({ playlist, active, onClick, onUpdate }) => {
  const [edit, setEdit] = useState(false)
  const [title, setTitle] = useState(playlist.name)
  const save = async () => {
    setEdit(false)
    try {
      playlist.name = title
      await updatePlaylist(playlist)
      onUpdate()
      toast.success("Playlist updating please wait");
    } catch (e) {
      toast.error("Unable to add video: invalid video url");
    }
  }

  const delPlaylist = async () => {
    setEdit(false)
    try {
      await deletePlaylist(playlist)
      onUpdate()
      toast.success("Playlist updating please wait");
    } catch (e) {
      toast.error("Unable to add video: invalid video url");
    }
  }


  return (
    <li onClick={onClick} className={`${active?"bg-violet-950":"bg-indigo-950"} flex items-center p-2 rounded-lg text-white hover:bg-violet-800`}>
      {(edit && active) ? <div className="relative w-full">
        <input type="text" value={title} onChange={(e) => setTitle(e.target.value)} className="block p-2.5 w-full z-20 text-sm text-gray-900 bg-gray-50 rounded-lg border-l-gray-100 border-l-2 border border-gray-300 focus:ring-violet-500 focus:border-violet-500 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:border-violet-500" required />
        <button onClick={() => save()} type="submit" className="absolute top-0 right-0 p-2.5 text-sm font-medium text-white bg-violet-700 rounded-r-lg border border-violet-700 hover:bg-violet-800 focus:ring-4 focus:outline-none focus:ring-violet-300 dark:bg-violet-600 dark:hover:bg-violet-700 dark:focus:ring-violet-800">
          <svg xmlns="http://www.w3.org/2000/svg" className="w-5 h-5" viewBox="0 0 24 24" strokeWidth="1.5" stroke="currentColor" fill="none" strokeLinecap="round" strokeLinejoin="round">
            <path stroke="none" d="M0 0h24v24H0z" fill="none"></path>
            <path d="M6 4h10l4 4v10a2 2 0 0 1 -2 2h-12a2 2 0 0 1 -2 -2v-12a2 2 0 0 1 2 -2"></path>
            <path d="M12 14m-2 0a2 2 0 1 0 4 0a2 2 0 1 0 -4 0"></path>
            <path d="M14 4l0 4l-6 0l0 -4"></path>
          </svg>
        </button>

      </div>
        :
        <span className="px-3 relative flex w-full justify-between items-center">
          <p className="py-2">{title}</p>
          {active &&
            <div className="gap-2 flex">
              <button type="button" onClick={() => setEdit(!edit)} className="text-white bg-purple-700 hover:bg-purple-800 focus:outline-none focus:ring-4 focus:ring-purple-300 font-medium rounded-full p-2 dark:bg-purple-600 dark:hover:bg-purple-700 dark:focus:ring-purple-900">
                <svg xmlns="http://www.w3.org/2000/svg" className="w-5 h-5" viewBox="0 0 24 24" strokeWidth="2" stroke="currentColor" fill="none" strokeLinecap="round" strokeLinejoin="round">
                  <path stroke="none" d="M0 0h24v24H0z" fill="none"></path>
                  <path d="M14 6l7 7l-4 4"></path>
                  <path d="M5.828 18.172a2.828 2.828 0 0 0 4 0l10.586 -10.586a2 2 0 0 0 0 -2.829l-1.171 -1.171a2 2 0 0 0 -2.829 0l-10.586 10.586a2.828 2.828 0 0 0 0 4z"></path>
                  <path d="M4 20l1.768 -1.768"></path>
                </svg>
              </button>

              <button type="button" onClick={() => delPlaylist()} className="text-white bg-red-700 hover:bg-red-800 focus:outline-none focus:ring-4 focus:ring-purple-300 font-medium rounded-full p-2 dark:bg-red-600 dark:hover:bg-red-700 dark:focus:ring-red-900">
                <svg fill="none" className="w-5 h-5" stroke="currentColor" strokeWidth={1.5} viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg" aria-hidden="true">
                  <path strokeLinecap="round" strokeLinejoin="round" d="M14.74 9l-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 01-2.244 2.077H8.084a2.25 2.25 0 01-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 00-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 013.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 00-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 00-7.5 0" />
                </svg>
              </button>
            </div>
          }
        </span>
      }
    </li>
  )
}

export const PlaylistBtn = ({ playlists }) => {
  const [show, setShow] = useState(false)

  const ref = useRef();
  useOnClickOutside(ref, () => setShow(false));

  const loadPlaylist = async (id) => {
    try {
      await loadFromPlaylist(id)
      
    } catch (e) {
      toast.error("Unable to load playlist");
    }

    setShow(false)
  }

  return <div className="w-full md:w-60 px-4 relative flex justify-end z-30" ref={ref}>
    <button onClick={() => setShow(!show)} className="rounded-full w-16 h-16 justify-center mb-2 text-white font-medium text-sm text-center inline-flex items-center bg-purple-600 hover:bg-violet-700 focus:ring-violet-800" type="button">

      <svg xmlns="http://www.w3.org/2000/svg" className="w-8 h-8" viewBox="0 0 24 24" strokeWidth="2" stroke="currentColor" fill="none" strokeLinecap="round" strokeLinejoin="round">
        <path stroke="none" d="M0 0h24v24H0z" fill="none"></path>
        <path d="M14 17m-3 0a3 3 0 1 0 6 0a3 3 0 1 0 -6 0"></path>
        <path d="M17 17v-13h4"></path>
        <path d="M13 5h-10"></path>
        <path d="M3 9l10 0"></path>
        <path d="M9 13h-6"></path>
      </svg>
      {/* <svg className="w-4 h-4 ml-2" aria-hidden="true" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"></path></svg> */}
    </button>

    <div className={`${!show && "hidden"} z-10 absolute bottom-20 w-60 right-4 rounded-lg bg-violet-950 shadow`}>
      {playlists.length > 0 && <ul className="max-h-60 py-2 overflow-y-auto hide-scrollbar text-violet-200" aria-labelledby="dropdownUsersButton">
        {playlists.map(playlist => (
          <li key={playlist.id} onClick={() => loadPlaylist(playlist.id)} className="flex cursor-pointer justify-center px-4 py-2 text-center  hover:bg-violet-600 hover:text-white">
            {playlist.name}

          </li>
        ))}
      </ul>}
      <Link to="playlists" className={`${playlists.length == 0 && "rounded-t-lg"} rounded-b-lg  flex justify-center p-3 text-sm font-medium text-white border-t border-violet-600  bg-violet-700 hover:bg-violet-600`}>
      manage playlists
      </Link>
    </div>
  </div>
}


const PlaylistPage = () => {
  const [playlists, setPlaylists] = useState([])
  const [playlist, setPlaylist] = useState({})

  useEffect(() => {
    const get = async () => {
      setPlaylists(await getChannelPlaylists())
    }
    get()
  }, [])

  const onUpdate = async () => {
    const _playlists = await getChannelPlaylists()
    setPlaylists(_playlists)
    if (playlist.id) {
      const _playlist = _playlists.filter(p => p.id === playlist.id)[0]
      if (_playlist !== undefined) {
        setPlaylist(_playlist)
      } else {
        setPlaylist({})
      }

    }
  }

  const newPlaylist = async () => {
    await createPlaylist()
    await onUpdate()
  }

  return (
    <>
      <main className='w-full flex flex-col md:flex-row  top-16 bottom-0'>
        <div className="w-full lg:w-1/4 md:w-1/2 h-64 md:h-full overflow-auto shadow-left text-white bg-zinc-900 hide-scrollbar">
          <div className="flex justify-between px-4 py-2.5  mb-2 sticky top-0 bg-zinc-800 z-10">
            <h1 className="text-2xl">Playlists</h1>
            <button onClick={() => newPlaylist()} className="text-white bg-purple-700 hover:bg-purple-800 focus:outline-none focus:ring-4 focus:ring-purple-300 font-medium rounded-full px-3 py-2 dark:bg-purple-600 dark:hover:bg-purple-700 dark:focus:ring-purple-900">
              New Playlist
            </button>
          </div>
          <ul className="space-y-2 font-medium px-2">
            {playlists.map(p => (
              <PlaylistHeader key={p.id} playlist={p} active={p.id === playlist.id} onUpdate={onUpdate} onClick={() => setPlaylist(p)} />
            ))}
          </ul>
        </div>
        {playlist.id && <ManagePlaylist playlist={playlist} onUpdate={onUpdate} />}
      </main>
    </>
  )
}
export default PlaylistPage