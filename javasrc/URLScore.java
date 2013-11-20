import java.net.MalformedURLException;
import java.net.URL;

// A data object for a URL and its score

public class URLScore implements Comparable<URLScore> {
    private static int Counter=0;
    private URL _url;
    private double _score;
    private int _index;

    public URLScore(String url, double score) throws MalformedURLException {
        _url = new URL(url);
        _score = score;
        _index = Counter++;
    }

    public URL getURL() {
        return _url;
    }

    public double getScore() {
        return _score;
    }

    public String getURLString() {
        return _url.toString();
    }

    public int getIndex() {
        return _index;
    }

    public int compareTo(URLScore other){
        if (other.getScore() > this.getScore())
            return 1;
        else if (other.getScore() < this.getScore())
            return -1;
            // else same score, then rely on the order in which the url was inserted into the queue
        else {
            if (other.getIndex() > this.getIndex())
                return -1;
            else
                return 1;
        }
    }
}

