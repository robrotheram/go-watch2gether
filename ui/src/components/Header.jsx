const Header = ({current}) => {
    return (
  <div className="flex md:items-end w-full shadow-head p-0 md:pt-8 md:pb-8 md:px-24 flex-col md:flex-row ">
    <img className="mr-0 md:rounded-xl shadow-xl w-full md:h-48 md:w-48 object-cover mt-0" src={current.thumbnail} />
    <div className="flex flex-col sm:ml-8 justify-start  md:justify-center p-4 md:p-0 md:pb-2">
      <h4 className="mt-0 mb-2 uppercase text-white tracking-widest text-xs">Now Playing</h4>
      <h1 className="mt-0 text-white text-3xl md:text-5xl">{current.title}</h1>
    </div>
  </div>
    )
}

export default Header