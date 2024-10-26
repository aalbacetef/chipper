import type { Status } from "./status";

const defaultDuration = 1500;

export type Notification = {
  id: string;
  text: string;
  status: Status;
  duration: number; // milliseconds
  at: Date;
}

export function NewNotification(id: string, status: Status, text: string, duration: number = defaultDuration): Notification {
  return {
    id,
    status,
    text,
    duration,
    at: new Date(),
  }
}
