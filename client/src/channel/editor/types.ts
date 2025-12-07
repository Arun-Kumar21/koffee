import { ReactElement } from "react";

export type RequestType = {
  userId: string;
  userName: string;
};

export interface ConnectUserData {
  userId : string;
  color : string;
  userName : string;
};

export interface UserCursor {
  user: {
    id: string;
    name: string;
    color: string;
  };
  cursor: {
    from: number;
    to: number;
  };
}

export type CommandType = {
  name : string;
  icon : ReactElement,
  action : () => boolean;
}

export type ConnectionStatusType = "connected" | "waiting" | "denied" | "error";
