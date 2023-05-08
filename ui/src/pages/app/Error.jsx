import { useSearchParams } from "react-router-dom";
import bg from '../../assets/404.jpg'

export function ErrroPage() {
    const [searchParams, setSearchParams] = useSearchParams();
    return (
      <div className="bg-indigo-900 relative overflow-hidden h-screen w-full">
        <img src={bg} className="absolute h-full w-full object-cover" />
        <div className="inset-0 bg-black opacity-25 absolute">
        </div>
        <div className="container mx-auto px-6 md:px-12 relative z-10 flex items-center py-32 xl:py-40">
          <div className="w-full font-mono flex flex-col items-center relative z-10">
            <h1 className="font-extrabold text-4xl text-center text-white leading-tight mt-4">
              {searchParams.length === 0 ? ("Please select a channel") : (<>Sorry channel is not registered. <br /><br /> Please connect the bot the the channel</>)}
            </h1>
          </div>
        </div>
      </div>
    )
  }
  