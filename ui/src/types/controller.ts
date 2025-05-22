export interface Media {
    id: string;
    user: string;
    url: string;
    audio_url?: string;
    type?: MediaType;
    title: string;
    channel?: string;
    time: MediaDuration;
    thumbnail?: string;
    order?: number;
    loading?: boolean;
}

export interface Playlist {
    id: string;
    name: string;
    videos: Media[];
    user: string;
    channel: string;
}

export interface MediaDuration {
    duration: number;
    progress: number;
}

export enum MediaType {
    YOUTUBE = 'YOUTUBE',
    YOUTUBE_LIVE = 'YOUTUBE_LIVE',
    SOUNDCLOUD = 'SOUNDCLOUD',
    ODYSEE = 'ODYSEE',
    MP4 = 'MP4',
    MP3 = 'MP3',
    PODCAST = 'PODCAST',
    RADIO_GARDEN = 'RADIO_GARDEN',
    PEERTUBE = 'PEERTUBE'
}

export enum PlayState {
    PLAY = 'PLAY',
    PAUSED = 'PAUSED',
    STOPPED = 'STOPPED'
}

export interface State {
    id: string;
    status: PlayState;
    queue: Media[];
    current: Media | null;
    loop: boolean;
    active: boolean;
}

export interface DiscordGuild {
    id: string;
    name: string;
    icon: string;
    owner: boolean;
    permissions: number;
    features: string[];
    permissions_new: string;
}

export type DiscordGuilds = DiscordGuild[];

export interface User {
    id: string;
    username: string;
    type: string;
    avatar: string;
    avatar_icon: string;
}

export interface Settings {
    bot: string;  // matches json:"bot" tag from Go struct
}

export interface Event {
    id: string;
    action: Action;
    state: State;
    message: string;
    players: PlayerMeta[];
}

type ActionType = 
    | "PLAY"
    | "PAUSE" 
    | "ADD_QUEUE"
    | "SEEK"
    | "UPDATE_QUEUE"
    | "UPDATE"
    | "REMOVE_QUEUE"
    | "UPDATE_DURATION"
    | "STOP"
    | "LOOP"
    | "SHUFFLE"
    | "SKIP"
    | "PlAYER_CHANGE"
    | "LEAVE";

export interface Action {
    type: ActionType;
    user: string;
    channel: string;
}

export interface PlayerMeta {
    id: string;
    user: string;
    progress: MediaDuration;
    type: string;
    running: boolean;
}