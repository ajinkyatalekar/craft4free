import { CButton } from "@/components/CButton";
import { useNavigate } from "react-router-dom";
import { Construction, ArrowLeft } from "lucide-react";

function UnderConstruction() {
  const navigate = useNavigate();

  const handleBackToHome = () => {
    navigate("/servers");
  };

  return (
    <>
      <div className="pt-20 bg-background flex items-center justify-center">
        <div className="rounded-xl p-8 shadow-lg max-w-md w-full">
          <div className="flex flex-col items-center text-center">
            <Construction className="w-20 h-20 text-yellow-100 mb-4" />

            <h1 className="text-2xl md:text-3xl font-bold mb-4">
              Under Construction
            </h1>

            <p className="text-muted-foreground mb-8">
              We're currently building this page. Check back soon!
            </p>

            <CButton
              onClick={handleBackToHome}
              className="w-full cursor-pointer"
            >
              <ArrowLeft className="mr-2 h-4 w-4" />
              Go Back
            </CButton>
          </div>
        </div>
      </div>
    </>
  );
}

export default UnderConstruction;
