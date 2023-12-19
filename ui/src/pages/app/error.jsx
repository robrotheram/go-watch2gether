import { Link, useSearchParams } from "react-router-dom";
// import bg from '../../assets/404.jpg'
import { useEffect, useState } from "react";
import { GuildIcon } from "./components/roomList";
import { getGuilds } from "./watch2gether";

export function ErrroPage() {
  const [guilds, setGuilds] = useState([]);

  useEffect(() => {
    async function get() {
      const g = await getGuilds();
      if (g != null) {
        setGuilds(g);
      }
    }
    if (guilds.length == 0) { get() };
  }, [])

  return (
    <div className={`text-white wrap-login min-h-screen w-full flex justify-center items-center bg-gradient-to-b from-violet-900  to-black`}>
      <div class="w-full max-w-4xl p-4  md:rounded-2xl shadow sm:p-6 bg-zinc-900 h-full md:h-3/4 flex flex-col z-10"  >
        {/* <h5 class="mb-3 font-semibold text-xl text-white">
          Channel List
        </h5> */}
       <ul class="my-4 space-y-3 flex-grow overflow-auto hide-scrollbar" >
        {
          guilds.map((guild) =>
            <li key={guild.id}>
              <Link key={guild.id} to={`/app/${guild.id}`} className="flex items-center p-3 text-base font-bold rounded-lg group hover:shadow bg-violet-900 hover:bg-violet-500 text-white">
                <GuildIcon guild={guild} />
                <span class="flex-1 ml-3 whitespace-nowrap">{guild.name}</span>
              </Link>
            </li>
          )
          }
        </ul>
      </div>


    </div>
  )
}