import org.jsoup.Jsoup;
import org.jsoup.select.Elements;
import org.jsoup.nodes.Element;
import org.jsoup.nodes.Document;

import java.io.IOException;
import java.net.URL;

/*
 * This parser uses the jsoup Java HTML Parser -- see http://jsoup.org/
 */
public class HTMLParser {
    private Document currentDoc;   // The last document retrieved

    public void connect(String url, String agent) throws IOException {
        try {
            currentDoc = Jsoup.connect(url).userAgent(agent).get();
        }
        catch (java.nio.charset.IllegalCharsetNameException e)  {
            throw new IOException(e.toString());
        }
    }

    public Elements getLinks() throws IOException {
	/********************************************************/
	/* GAP!							*/
	/* Get the links from the last document retrieved	*/
	/* We are only interested in <a href> links that	*/
	/* are html pages.					*/
	/********************************************************/
    }

    public String getBody() throws IOException {
	/********************************************************/
	/* GAP!							*/
	/* Get the text of the body from the last document 	*/
	/* retrieved						*/
	/********************************************************/ 
    }
}
