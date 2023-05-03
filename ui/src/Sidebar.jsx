import RoomList from './RoomList'

const Sidebar = ({guilds}) => {
    return <div className='h-full p-0.5 hidden md:block overflow-y-auto hide-scrollbar' >
       <RoomList guilds={guilds}/>
    </div>
}

export default Sidebar