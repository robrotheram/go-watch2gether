const Header = ({current}) => {
    return (
  <div className="flex items-end w-full shadow-head p-0 md:pt-24 md:pb-16 md:px-24 flex-col md:flex-row ">
    <img className="mr-0 md:mr-6  md:rounded-xl shadow-xl w-full md:h-48 md:w-48 object-cover  " src={current.thumbnail} />
    <div className="flex flex-col justify-center p-4 md:p-0">
      <h4 className="mt-0 mb-2 uppercase text-white tracking-widest text-xs">Now Playing</h4>
      <h1 className="mt-0 text-white text-3xl md:text-5xl">{current.title}</h1>
    </div>
  </div>
    )
}

export default Header