import { CButton } from "@/components/CButton";
import { Logo } from "@/components/SiteLogo";
import { useNavigate } from "react-router-dom";

function Landing() {
  const navigate = useNavigate();
  const handleAuthRedirect = () => {
    navigate("/login");
  };

  return (
    <>
      <div className="h-screen bg-gray-800 bg-[url(assets/background/caves.png)] bg-blend-overlay bg-cover flex flex-col">
        <div className="flex flex-col items-center justify-center flex-grow gap-4 px-20">
          <Logo otherStyles="w-[10rem] sm:w-[20rem]" />
          <p className="pt-5" />
          <p className="text-5xl text-center sm:text-7xl">
            Free Minecraft Servers
          </p>
          <p className="text-xl text-center text-muted-foreground">
            Play Minecraft with your friends for free!
          </p>
          <div className="mt-2" />
          <CButton
            onClick={() => {
              handleAuthRedirect();
            }}
          >
            Play
          </CButton>
        </div>
      </div>

      <div className="py-20 bg-gray-900">
        <div className="container mx-auto px-4">
          <div className="text-center mb-16">
            <h2 className="text-4xl mb-4">How We Keep Servers Free</h2>
            <p className="text-xl text-muted-foreground max-w-2xl mx-auto">
              Craft4Free provides premium Minecraft servers at zero cost to our
              players.
            </p>
          </div>

          <div className="grid grid-cols-1 lg:grid-cols-2 gap-12">
            <div className="bg-background rounded-xl p-8 shadow-lg flex flex-col justify-between h-full">
              <div>
                <h3 className="text-3xl mb-6">Powered by Community</h3>
                <p className="text-lg mb-8">
                  We believe Minecraft should be accessible to everyone. That's
                  why we host and maintain high-performance servers completely
                  free of charge.
                </p>

                <div className="flex items-center mb-6">
                  <div className="bg-green-500 p-3 rounded-full mr-4">
                    <svg
                      xmlns="http://www.w3.org/2000/svg"
                      className="h-6 w-6 text-white"
                      fill="none"
                      viewBox="0 0 24 24"
                      stroke="currentColor"
                    >
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth={2}
                        d="M5 13l4 4L19 7"
                      />
                    </svg>
                  </div>
                  <span className="text-lg">Premium Hardware</span>
                </div>

                <div className="flex items-center mb-6">
                  <div className="bg-green-500 p-3 rounded-full mr-4">
                    <svg
                      xmlns="http://www.w3.org/2000/svg"
                      className="h-6 w-6 text-white"
                      fill="none"
                      viewBox="0 0 24 24"
                      stroke="currentColor"
                    >
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth={2}
                        d="M5 13l4 4L19 7"
                      />
                    </svg>
                  </div>
                  <span className="text-lg">24/7 Uptime</span>
                </div>

                <div className="flex items-center">
                  <div className="bg-green-500 p-3 rounded-full mr-4">
                    <svg
                      xmlns="http://www.w3.org/2000/svg"
                      className="h-6 w-6 text-white"
                      fill="none"
                      viewBox="0 0 24 24"
                      stroke="currentColor"
                    >
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth={2}
                        d="M5 13l4 4L19 7"
                      />
                    </svg>
                  </div>
                  <span className="text-lg">Instant Setup</span>
                </div>
              </div>
            </div>

            <div className="bg-background rounded-xl p-8 shadow-lg flex flex-col justify-between h-full">
              <div>
                <h3 className="text-3xl mb-6">Supported by Ads</h3>
                <p className="text-lg mb-6">
                  Instead of charging our players, we cover server costs through
                  minimally invasive advertisements.
                </p>
                <p className="text-lg mb-8">
                  This approach allows us to provide a premium gaming experience
                  while keeping it 100% free for everyone to enjoy.
                </p>
                <p className="text-lg font mb-8">
                  Please disable ad blocking or whitelist craft4free to support
                  our services and keep craft4free free :)
                </p>
              </div>
              <div>
                <CButton
                  onClick={() => {
                    handleAuthRedirect();
                  }}
                >
                  Create Your Server
                </CButton>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div className="bg-background p-5">
        <p className="text-sm text-center">© 2025 by Ajinkya Talekar</p>
      </div>
    </>
  );
}

export default Landing;
