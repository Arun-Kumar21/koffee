import { useState, useEffect } from "react";
import { Button } from "../ui/button";
import { ToastAction } from "../ui/toast";
import { useToast } from "../../hooks/use-toast";

interface ConnectAccessToastProps {
  userName : string;
  onAccept: () => void;
  onReject: () => void;
}

export default function ConnectAccessToast({
  userName,
  onAccept,
  onReject,
}: ConnectAccessToastProps) {
  const { toast } = useToast();
  const [isShown, setIsShown] = useState(false);

  useEffect(() => {
    if (!isShown) {
      setIsShown(true);
      toast({
        title: "Connect Request",
        description:  `${userName} wants to connect with the channel.`,
        action: (
          <div className="flex space-x-2">
            <ToastAction altText="Reject" onClick={handleReject}>
              Reject
            </ToastAction>
            <Button variant="default" size="sm" onClick={handleAccept}>
              Accept
            </Button>
          </div>
        ),
        duration: Infinity,
      });
    }
  }, [isShown]);

  const handleAccept = () => {
    setIsShown(false);
    onAccept();
  };

  const handleReject = () => {
    setIsShown(false);
    onReject();
  };

  return null;
}
