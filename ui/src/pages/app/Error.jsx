import { Link, useSearchParams } from "react-router-dom";
import bg from '../../assets/404.jpg'
import { useEffect, useState } from "react";
import { getGuilds } from "../../api/watch2gether";
import { GuildIcon } from "../../components/RoomList";

export function ErrroPage() {
  const [searchParams, setSearchParams] = useSearchParams();
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
    <div className="bg-indigo-900 relative overflow-hidden h-screen w-full flex flex-col justify-center items-center">
      <img src={bg} className="absolute h-full w-full object-cover " />
      <div class="w-full max-w-4xl p-4  md:rounded-2xl shadow sm:p-6 bg-zinc-900 h-full md:h-1/2 flex flex-col z-10"  >
        <h5 class="mb-3 font-semibold text-gray-900 text-xl dark:text-white">
          Channel List
        </h5>
       <ul class="my-4 space-y-3 flex-grow overflow-auto hide-scrollbar" >
        {
          guilds.map((guild) =>
            <li key={guild.id}>
              <Link key={guild.id} to={`/app/${guild.id}`} className="flex items-center p-3 text-base font-bold text-gray-950 rounded-lg bg-gray-50 hover:bg-gray-100 group hover:shadow dark:bg-violet-900 dark:hover:bg-violet-500 dark:text-white">
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
