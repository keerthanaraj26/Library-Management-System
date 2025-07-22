async function addbook(event) {
  event.preventDefault(); 

  const form = document.getElementById("bookform");
  const formData = new FormData(form);

  const data = {
    id: parseInt(formData.get("id")) || 0,
    title: formData.get("title"),
    author: formData.get("author"),
    count: parseInt(formData.get("count")) || 0
  };

  try {
    const response = await fetch("/api/addbook", {
      method: "POST",
      headers: {
        "Content-Type": "application/json"
      },
      body: JSON.stringify(data)
    });

    const result = await response.json();
    console.log("Response:", result);
    window.location.href = "/viewcatalogue";

    } catch (error) {
      console.error("Error:", error);
      document.getElementById("responseMsg").innerText = "Submission failed!";
    }
}
