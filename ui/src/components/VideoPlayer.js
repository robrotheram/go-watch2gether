import React from "react"
import ReactPlayer from 'react-player'
import {connect} from 'react-redux'
import {play, pause,updateSeek, updateQueue} from '../store/room/room.actions'
import { Empty } from "antd"

class VideoPlayer extends React.Component {
    state = {
        url: null,
        pip: false,
        playing: true,
        light: false,
        volume: 0.8,
        muted: false,
        played: 0,
        loaded: 0,
        duration: 0,
        playbackRate: 1.0,
        loop: false
    }


    componentWillReceiveProps(nextProps) {
        // You don't have to do this check first, but it can help prevent an unneeded render
        if (nextProps.seek !== this.state.played) {
            console.log("STATE_CHANGE",nextProps )
          this.setState({ played: nextProps.seek });
          if (this.player !== undefined) {
            this.player.seekTo(parseFloat(nextProps.seek))
          }
        }
      }

    componentDidMount() {
        let _this = this
        setTimeout(function(){
            if (_this.player !== undefined && _this.player !== null){
                _this.player.seekTo(parseFloat(_this.props.seek)) 
            }
        }, 200);
        
    }

    handlePlay = () => {
        if(!this.props.playing){
            play();
        }
        
    }

    handlePause = () => {
        if(this.props.playing){
            pause();
        }
    }

    handleEnded = () => {
        let videoList = [...this.props.queue];
        videoList.splice(0, 1);
        let lowestSeek = 1; 
        this.props.users.forEach(user => {
            if (lowestSeek > user.seek){
                lowestSeek = user.seek
            }
        });
        if (lowestSeek >= 0.9) {
            updateQueue(videoList)
        }
    }

    handleProgress = state => {
        // console.log('onProgress', state)
        // console.log("PLAYER!", this.player)
        this.setState({ seek: state.playedSeconds })
        updateSeek(state.played)
    }

    onProgressUpdate = e => {
        this.setState({ seeking: false })
        this.player.seekTo(parseFloat(e.target.value))
    }

    
    ref = player => { this.player = player }

    render(){
            const {queue, playing } = this.props
            return(
            <div style={{ "height":"500px", "width":"100%"}}>  
            {queue[0] !== undefined ? 
               <ReactPlayer 
                ref={this.ref}
                width="100%" height="500px"  
                url={queue[0].url} 
                controls={true}
                playing={playing}
                onPause={this.handlePause}
                onPlay={this.handlePlay}
                onProgress={this.handleProgress}
                onEnded={this.handleEnded}
                />  
            : <Empty  style={{ "height":"500px", "width":"100%", "paddingTop":"180px"}}/>}
          </div>  
        )
    }
}
const mapStateToProps  = (state) =>{
    return state.room
  } 
export default connect(mapStateToProps, {updateQueue, updateSeek})(VideoPlayer)
  