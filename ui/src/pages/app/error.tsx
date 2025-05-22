import { Link } from "react-router-dom";
import { GuildIcon } from "./components/roomList";
import { useContext } from "react";
import { GuildContext } from "./components/providers";

export function ErrroPage() {
  const { guilds } = useContext(GuildContext)
  return (
    <div className='text-white wrap-login min-h-screen w-full flex justify-center items-center'>
      <div className="w-full max-w-4xl p-4  md:rounded-2xl shadow sm:p-6 bg-zinc-900 h-full md:h-3/4 flex flex-col z-10"  >
        <ul className="my-4 space-y-3 flex-grow overflow-auto hide-scrollbar" >
          {
            guilds.map((guild) =>
              <li key={guild.id}>
                <Link key={guild.id} to={`/app/${guild.id}`} className="flex items-center p-3 text-base font-bold rounded-lg group hover:shadow bg-violet-900 hover:bg-violet-500 text-white">
                  <GuildIcon guild={guild} />
                  <span className="flex-1 ml-3 whitespace-nowrap">{guild.name}</span>
                </Link>
              </li>
            )
          }
        </ul>
      </div>
    </div>
  )
}